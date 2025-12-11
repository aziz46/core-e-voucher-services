package handler

import (
	"context"
	"log"

	"github.com/aziz46/core-e-voucher-services/pkg/db"
	"github.com/aziz46/core-e-voucher-services/pkg/models"
	"github.com/gofiber/fiber/v2"
)

// ReserveRequest represents a request to reserve credit limit
type ReserveRequest struct {
	TransactionID string `json:"tx_id"`
	Amount        int64  `json:"amount"`
}

// RestoreRequest represents a request to restore credit limit
type RestoreRequest struct {
	TransactionID string `json:"tx_id"`
	Amount        int64  `json:"amount"`
}

// LimitResponse represents the current credit limit
type LimitResponse struct {
	PartnerID      string `json:"partner_id"`
	LimitTotal     int64  `json:"limit_total"`
	LimitUsed      int64  `json:"limit_used"`
	LimitAvailable int64  `json:"limit_available"`
}

// GetLimit returns the current credit limit for a partner
func GetLimit(c *fiber.Ctx) error {
	partnerID := c.Params("partner_id")
	ctx := context.Background()

	row := db.Pool.QueryRow(ctx,
		"SELECT partner_id, limit_total, limit_used, limit_available FROM credit_limits WHERE partner_id = $1",
		partnerID)

	var limit models.CreditLimit
	err := row.Scan(&limit.PartnerID, &limit.LimitTotal, &limit.LimitUsed, &limit.LimitAvailable)
	if err != nil {
		log.Printf("Error fetching limit: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "limit not found",
		})
	}

	return c.JSON(LimitResponse{
		PartnerID:      limit.PartnerID,
		LimitTotal:     limit.LimitTotal,
		LimitUsed:      limit.LimitUsed,
		LimitAvailable: limit.LimitAvailable,
	})
}

// Reserve atomically reduces the available limit
func Reserve(c *fiber.Ctx) error {
	partnerID := c.Params("partner_id")
	var req ReserveRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount must be greater than 0",
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

	// Use SELECT FOR UPDATE to lock the row
	var limitAvailable int64
	err = tx.QueryRow(ctx,
		"SELECT limit_available FROM credit_limits WHERE partner_id = $1 FOR UPDATE",
		partnerID).Scan(&limitAvailable)

	if err != nil {
		log.Printf("Error fetching limit: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "limit not found",
		})
	}

	// Check if enough balance
	if limitAvailable < req.Amount {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":     "insufficient credit limit",
			"available": limitAvailable,
			"requested": req.Amount,
		})
	}

	// Update limit
	_, err = tx.Exec(ctx,
		"UPDATE credit_limits SET limit_used = limit_used + $1, limit_available = limit_available - $1 WHERE partner_id = $2",
		req.Amount, partnerID)

	if err != nil {
		log.Printf("Error updating limit: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to reserve limit",
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "transaction commit failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "reserved",
		"amount": req.Amount,
	})
}

// Restore returns the reserved amount back to available limit
func Restore(c *fiber.Ctx) error {
	partnerID := c.Params("partner_id")
	var req RestoreRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount must be greater than 0",
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

	// Use SELECT FOR UPDATE
	var limitUsed int64
	err = tx.QueryRow(ctx,
		"SELECT limit_used FROM credit_limits WHERE partner_id = $1 FOR UPDATE",
		partnerID).Scan(&limitUsed)

	if err != nil {
		log.Printf("Error fetching limit: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "limit not found",
		})
	}

	// Check if enough used limit to restore
	if limitUsed < req.Amount {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":     "insufficient used limit to restore",
			"used":      limitUsed,
			"requested": req.Amount,
		})
	}

	// Update limit
	_, err = tx.Exec(ctx,
		"UPDATE credit_limits SET limit_used = limit_used - $1, limit_available = limit_available + $1 WHERE partner_id = $2",
		req.Amount, partnerID)

	if err != nil {
		log.Printf("Error updating limit: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to restore limit",
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "transaction commit failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "restored",
		"amount": req.Amount,
	})
}
