package models

import "time"

// Tenant represents a tenant in the system
type Tenant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Plan      string    `json:"plan" db:"plan"`
	APIKey    string    `json:"-" db:"api_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Partner represents a partner/reseller
type Partner struct {
	ID        string    `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Name      string    `json:"name" db:"name"`
	KYCStatus string    `json:"kyc_status" db:"kyc_status"`
	Contact   string    `json:"contact" db:"contact"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreditLimit represents credit limit for a partner
type CreditLimit struct {
	PartnerID      string    `json:"partner_id" db:"partner_id"`
	LimitTotal     int64     `json:"limit_total" db:"limit_total"`
	LimitUsed      int64     `json:"limit_used" db:"limit_used"`
	LimitAvailable int64     `json:"limit_available" db:"limit_available"`
	ResetPeriod    string    `json:"reset_period" db:"reset_period"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Transaction represents a PPOB transaction
type Transaction struct {
	ID             string    `json:"id" db:"id"`
	TenantID       string    `json:"tenant_id" db:"tenant_id"`
	PartnerID      string    `json:"partner_id" db:"partner_id"`
	Amount         int64     `json:"amount" db:"amount"`
	Fee            int64     `json:"fee" db:"fee"`
	Total          int64     `json:"total" db:"total"`
	Status         string    `json:"status" db:"status"`
	Provider       string    `json:"provider" db:"provider"`
	ProviderTxID   string    `json:"provider_tx_id" db:"provider_tx_id"`
	IdempotencyKey string    `json:"-" db:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Invoice represents a billing invoice
type Invoice struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	PartnerID   string    `json:"partner_id" db:"partner_id"`
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time `json:"period_end" db:"period_end"`
	AmountDue   int64     `json:"amount_due" db:"amount_due"`
	Status      string    `json:"status" db:"status"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           string    `json:"id" db:"id"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	ResourceID   string    `json:"resource_id" db:"resource_id"`
	Action       string    `json:"action" db:"action"`
	PayloadJSON  string    `json:"payload_json" db:"payload_json"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Transaction status constants
const (
	TxStatusPending   = "pending"
	TxStatusSuccess   = "success"
	TxStatusFailed    = "failed"
	TxStatusCancelled = "cancelled"
)

// Invoice status constants
const (
	InvoiceStatusDraft     = "draft"
	InvoiceStatusIssued    = "issued"
	InvoicStatusPaid       = "paid"
	InvoiceStatusOverdue   = "overdue"
	InvoiceStatusCancelled = "cancelled"
)
