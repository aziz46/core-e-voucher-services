package handler

import (
	"context"
	"log"
	"time"

	"github.com/aziz46/core-e-voucher-services/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ListInvoicesRequest is the query for listing invoices
type ListInvoicesRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// GenerateInvoiceRequest is the request to generate invoice
type GenerateInvoiceRequest struct {
	Period string `json:"period"` // "daily" or "monthly"
}

// ListInvoices returns all invoices for a tenant
func ListInvoices(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	ctx := context.Background()
	rows, err := db.Pool.Query(ctx,
		"SELECT id, tenant_id, partner_id, period_start, period_end, amount_due, status, due_date, created_at FROM invoices WHERE tenant_id = $1 LIMIT $2 OFFSET $3",
		tenantID, limit, offset)

	if err != nil {
		log.Printf("Error querying invoices: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to list invoices",
		})
	}
	defer rows.Close()

	type Invoice struct {
		ID          string    `json:"id"`
		TenantID    string    `json:"tenant_id"`
		PartnerID   string    `json:"partner_id"`
		PeriodStart time.Time `json:"period_start"`
		PeriodEnd   time.Time `json:"period_end"`
		AmountDue   int64     `json:"amount_due"`
		Status      string    `json:"status"`
		DueDate     time.Time `json:"due_date"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		err := rows.Scan(&inv.ID, &inv.TenantID, &inv.PartnerID, &inv.PeriodStart, &inv.PeriodEnd,
			&inv.AmountDue, &inv.Status, &inv.DueDate, &inv.CreatedAt)
		if err != nil {
			log.Printf("Error scanning invoice: %v", err)
			continue
		}
		invoices = append(invoices, inv)
	}

	return c.JSON(fiber.Map{
		"invoices": invoices,
		"count":    len(invoices),
	})
}

// GenerateInvoice generates invoices for active partners
func GenerateInvoice(c *fiber.Ctx) error {
	tenantID := c.Params("tenant_id")
	var req GenerateInvoiceRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	defer tx.Rollback(ctx)

	// Get all transactions for tenant in current month
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	periodEnd := periodStart.AddDate(0, 1, 0).Add(-time.Second)

	rows, err := tx.Query(ctx,
		`SELECT DISTINCT partner_id, SUM(total) as total_amount 
		 FROM transactions 
		 WHERE tenant_id = $1 AND status = 'success' AND created_at BETWEEN $2 AND $3 
		 GROUP BY partner_id`,
		tenantID, periodStart, periodEnd)

	if err != nil {
		log.Printf("Error querying transactions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate invoices",
		})
	}
	defer rows.Close()

	var generated int
	for rows.Next() {
		var partnerID string
		var totalAmount int64
		if err := rows.Scan(&partnerID, &totalAmount); err != nil {
			continue
		}

		invoiceID := uuid.New().String()
		dueDate := periodEnd.AddDate(0, 0, 30) // 30 days due

		_, err := tx.Exec(ctx,
			`INSERT INTO invoices (id, tenant_id, partner_id, period_start, period_end, amount_due, status, due_date, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			invoiceID, tenantID, partnerID, periodStart, periodEnd, totalAmount, "issued", dueDate, now)

		if err != nil {
			log.Printf("Error creating invoice: %v", err)
			continue
		}
		generated++
	}

	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "transaction commit failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "invoices generated",
		"count":   generated,
	})
}

// PaymentCallback handles payment callbacks from providers
func PaymentCallback(c *fiber.Ctx) error {
	type CallbackRequest struct {
		TenantID  string `json:"tenant_id"`
		PartnerID string `json:"partner_id"`
		InvoiceID string `json:"invoice_id"`
		Amount    int64  `json:"amount"`
	}

	var req CallbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	defer tx.Rollback(ctx)

	// Update invoice status to paid
	_, err = tx.Exec(ctx,
		"UPDATE invoices SET status = $1 WHERE id = $2",
		"paid", req.InvoiceID)

	if err != nil {
		log.Printf("Error updating invoice: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to process payment",
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "transaction commit failed",
		})
	}

	return c.JSON(fiber.Map{
		"status": "payment processed",
	})
}
