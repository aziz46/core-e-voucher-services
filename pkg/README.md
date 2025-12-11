# PKG Directory

This directory contains shared libraries and utilities that can be used across all services.

## Structure

- `logger/` - Structured logging utilities (zap/zerolog)
- `config/` - Configuration management (viper/envconfig)
- `database/` - Database connection and utilities
- `redis/` - Redis client and utilities
- `middleware/` - HTTP middleware (auth, logging, rate limiting, etc.)

These packages are designed to be reusable and service-agnostic.
