-- Migration: 001_init.sql
-- Description: Create initial database schema for e-voucher PPOB system

-- Create tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    plan VARCHAR(50) NOT NULL DEFAULT 'basic',
    api_key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create partners table
CREATE TABLE IF NOT EXISTS partners (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    kyc_status VARCHAR(50) DEFAULT 'pending',
    contact VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Create credit_limits table
CREATE TABLE IF NOT EXISTS credit_limits (
    partner_id VARCHAR(36) PRIMARY KEY,
    limit_total BIGINT NOT NULL DEFAULT 1000000,
    limit_used BIGINT NOT NULL DEFAULT 0,
    limit_available BIGINT NOT NULL DEFAULT 1000000,
    reset_period VARCHAR(50) DEFAULT 'daily',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (partner_id) REFERENCES partners(id)
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    partner_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    fee BIGINT NOT NULL DEFAULT 0,
    total BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    provider VARCHAR(100) DEFAULT 'mock_provider',
    provider_tx_id VARCHAR(255),
    idempotency_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id)
);

-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    partner_id VARCHAR(36) NOT NULL,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    amount_due BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft',
    due_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id)
);

-- Create receivables table
CREATE TABLE IF NOT EXISTS receivables (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    partner_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'outstanding',
    tx_id VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id),
    FOREIGN KEY (tx_id) REFERENCES transactions(id)
);

-- Create provider_configs table
CREATE TABLE IF NOT EXISTS provider_configs (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    provider_name VARCHAR(100) NOT NULL,
    creds_encrypted TEXT,
    endpoint VARCHAR(255),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    resource_type VARCHAR(100),
    resource_id VARCHAR(36),
    action VARCHAR(100),
    payload_json TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_partners_tenant_id ON partners(tenant_id);
CREATE INDEX IF NOT EXISTS idx_credit_limits_partner_id ON credit_limits(partner_id);
CREATE INDEX IF NOT EXISTS idx_transactions_tenant_id ON transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_partner_id ON transactions(partner_id);
CREATE INDEX IF NOT EXISTS idx_transactions_idempotency_key ON transactions(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_invoices_tenant_id ON invoices(tenant_id);
CREATE INDEX IF NOT EXISTS idx_invoices_partner_id ON invoices(partner_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_receivables_tenant_id ON receivables(tenant_id);
CREATE INDEX IF NOT EXISTS idx_receivables_partner_id ON receivables(partner_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
