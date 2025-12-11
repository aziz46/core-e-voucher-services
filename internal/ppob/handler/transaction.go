package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aziz46/core-e-voucher-services/pkg/connector"
	"github.com/aziz46/core-e-voucher-services/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateTransactionRequest is the request to create a transaction
type CreateTransactionRequest struct {
	ProductCode    string `json:"product_code"`
	CustomerNo     string `json:"customer_no"`
	Amount         int64  `json:"amount"`
	PartnerID      string `json:"partner_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

// TransactionResponse is the response for a transaction
type TransactionResponse struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	Amount       int64     `json:"amount"`
	Fee          int64     `json:"fee"`
	Total        int64     `json:"total"`
	ProviderTxID string    `json:"provider_tx_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateTransaction creates a new PPOB transaction
func CreateTransaction(c *fiber.Ctx) error {
	tenantID := c.Params("tenant")
	var req CreateTransactionRequest

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

	// Check idempotency
	if req.IdempotencyKey != "" {
		var existingTxID string
		err := db.Pool.QueryRow(ctx,
			"SELECT id FROM transactions WHERE idempotency_key = $1 AND tenant_id = $2 LIMIT 1",
			req.IdempotencyKey, tenantID).Scan(&existingTxID)

		if err == nil {
			// Transaction already exists, return existing
			return getTransaction(c, tenantID, existingTxID)
		}
	}

	// Calculate fee (2.5% for demo)
	fee := req.Amount / 40
	total := req.Amount + fee
	txID := uuid.New().String()

	// Start database transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}
	defer tx.Rollback(ctx)

	// Create transaction record
	now := time.Now()
	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (id, tenant_id, partner_id, amount, fee, total, status, provider, idempotency_key, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		txID, tenantID, req.PartnerID, req.Amount, fee, total, "pending", "mock_provider", req.IdempotencyKey, now, now)

	if err != nil {
		log.Printf("Error creating transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create transaction",
		})
	}

	// Reserve credit from credit-service
	reserveReq := map[string]interface{}{
		"tx_id":  txID,
		"amount": total,
	}

	reserveBody, _ := json.Marshal(reserveReq)
	creditURL := fmt.Sprintf("http://credit-service:8080/v1/partners/%s/reserve", req.PartnerID)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", creditURL, bytes.NewBuffer(reserveBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Service-Token", "internal-service-token")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(httpReq)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Failed to reserve credit: %v", err)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "insufficient credit limit",
		})
	}
	resp.Body.Close()

	// Call provider payment
	provider := connector.NewMockProvider(0) // 0% failure rate for demo
	payReq := connector.PayRequest{
		ProductCode: req.ProductCode,
		CustomerNo:  req.CustomerNo,
		Amount:      req.Amount,
		RefNo:       txID,
	}

	payResp, err := provider.Pay(ctx, payReq)

	if err != nil {
		// Restore credit on failure
		restoreReq := map[string]interface{}{
			"tx_id":  txID,
			"amount": total,
		}
		restoreBody, _ := json.Marshal(restoreReq)
		restoreURL := fmt.Sprintf("http://credit-service:8080/v1/partners/%s/restore", req.PartnerID)

		httpReq, _ := http.NewRequestWithContext(ctx, "POST", restoreURL, bytes.NewBuffer(restoreBody))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Service-Token", "internal-service-token")
		client.Do(httpReq)

		// Mark transaction as failed
		_, _ = tx.Exec(ctx, "UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3",
			"failed", now, txID)
		tx.Commit(ctx)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payment failed",
		})
	}

	// Update transaction with provider response
	_, err = tx.Exec(ctx,
		"UPDATE transactions SET status = $1, provider_tx_id = $2, updated_at = $3 WHERE id = $4",
		"success", payResp.ProviderRefNo, now, txID)

	if err != nil {
		log.Printf("Error updating transaction: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "transaction commit failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(TransactionResponse{
		ID:           txID,
		Status:       "success",
		Amount:       req.Amount,
		Fee:          fee,
		Total:        total,
		ProviderTxID: payResp.ProviderRefNo,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

// GetTransaction retrieves a transaction by ID
func GetTransaction(c *fiber.Ctx) error {
	tenantID := c.Params("tenant")
	txID := c.Params("tx_id")

	return getTransaction(c, tenantID, txID)
}

func getTransaction(c *fiber.Ctx, tenantID, txID string) error {
	ctx := context.Background()

	row := db.Pool.QueryRow(ctx,
		`SELECT id, status, amount, fee, total, provider_tx_id, created_at, updated_at 
		 FROM transactions WHERE id = $1 AND tenant_id = $2`,
		txID, tenantID)

	var tx struct {
		ID           string
		Status       string
		Amount       int64
		Fee          int64
		Total        int64
		ProviderTxID string
		CreatedAt    time.Time
		UpdatedAt    time.Time
	}

	err := row.Scan(&tx.ID, &tx.Status, &tx.Amount, &tx.Fee, &tx.Total, &tx.ProviderTxID, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		log.Printf("Error fetching transaction: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transaction not found",
		})
	}

	return c.JSON(TransactionResponse{
		ID:           tx.ID,
		Status:       tx.Status,
		Amount:       tx.Amount,
		Fee:          tx.Fee,
		Total:        tx.Total,
		ProviderTxID: tx.ProviderTxID,
		CreatedAt:    tx.CreatedAt,
		UpdatedAt:    tx.UpdatedAt,
	})
}
