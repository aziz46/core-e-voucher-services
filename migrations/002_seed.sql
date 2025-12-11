-- Migration: 002_seed.sql
-- Description: Insert seed data for development and testing

-- Insert test tenant
INSERT INTO tenants (id, name, plan, api_key, created_at) 
VALUES ('tenant_001', 'PT Mitra Indonesia', 'premium', 'sk_live_abc123xyz789', NOW())
ON CONFLICT DO NOTHING;

-- Insert test partner
INSERT INTO partners (id, tenant_id, name, kyc_status, contact, created_at)
VALUES ('partner_001', 'tenant_001', 'Toko Elektronik Jaya', 'verified', '0812-3456-7890', NOW())
ON CONFLICT DO NOTHING;

-- Insert credit limit
INSERT INTO credit_limits (partner_id, limit_total, limit_used, limit_available, reset_period, updated_at)
VALUES ('partner_001', 1000000, 0, 1000000, 'daily', NOW())
ON CONFLICT DO NOTHING;

-- Insert provider config
INSERT INTO provider_configs (id, tenant_id, provider_name, endpoint, active, created_at)
VALUES ('provider_config_001', 'tenant_001', 'mock_provider', 'http://mock-provider:8080', true, NOW())
ON CONFLICT DO NOTHING;
