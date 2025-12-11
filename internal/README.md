# Internal Directory

This directory contains the business logic and internal packages for each service.

## Structure

- `credit-service/` - Credit service business logic (handlers, services, repositories)
- `billing-service/` - Billing service business logic (handlers, services, workers)
- `ppob-core/` - PPOB core service business logic (handlers, services, connectors)
- `common/` - Shared internal packages used across services

Each service directory typically contains:
- `handler/` - HTTP handlers
- `service/` - Business logic
- `repository/` - Data access layer
- `model/` - Domain models
