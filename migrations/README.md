# Migrations Directory

This directory contains SQL migration files for database schema management.

## File Naming Convention

Migrations are named with the following pattern:
```
{version}_{description}.sql
```

Example:
- `001_init.sql` - Initial schema setup
- `002_add_audit_logs.sql` - Add audit logging tables

## Tables

The initial migration creates the following tables:
- `tenants` - Tenant/client information
- `partners` - Partner/agent information
- `credit_limits` - Credit limit tracking per partner
- `transactions` - Transaction records
- `invoices` - Billing invoices
- `receivables` - Receivable records
- `provider_configs` - Provider configuration and credentials
- `audit_logs` - Audit trail for all operations

## Running Migrations

Migrations are automatically applied when services start up, or can be run manually using migration tools.
