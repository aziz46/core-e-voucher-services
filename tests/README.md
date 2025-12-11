# Tests Directory

This directory contains integration tests and load tests for the services.

## Structure

- `integration/` - Integration tests that test multiple services working together
- `load/` - Load and performance tests

## Integration Tests

Integration tests validate end-to-end flows:
- Transaction creation flow (credit check → provider payment → invoice creation)
- Credit limit reserve and restore operations
- Payment callback handling
- Idempotency verification

These tests require:
- Running instances of all services
- Database with test data
- Mock or test provider endpoints

## Load Tests

Load tests verify system performance under load:
- Transaction throughput (target: 500 TPS per instance)
- Latency measurements (target: P95 < 300ms)
- Concurrent credit operations
- Provider timeout handling

Load tests can be run using tools like:
- k6
- Apache JMeter
- Vegeta
