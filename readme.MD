# Instruksi untuk Agent Copilot — Implementasi PaaS PPOB (Postpaid / No-Deposit)

**Tujuan:** Agent Copilot membuatkan kode & artefak untuk MVP PaaS PPOB model postpaid (mitra/teller bisa transaksi tanpa deposit) — fokus ke 3 service utama: **Credit-Limit Service**, **Billing Service**, dan **Core PPOB Service + Provider Connector**.

> Dokumentasi ini ditulis supaya agent bisa generate implementasi end-to-end (kode, config, infra dev, test). Baca seluruh dokumen sebelum mulai.

---

## 1. Deliverables yang harus dibuat oleh agent

1. Repo skeleton (`/core-e-voucher-services`) dengan subfolders:

   * `cmd/` (main apps)
   * `internal/` (business logic)
   * `pkg/` (shared libs)
   * `deploy/` (docker-compose + k8s manifests)
   * `migrations/` (sql)
   * `docs/` (OpenAPI + README)
   * `tests/` (integration + load test)
2. Implementasi 3 layanan (Golang, gin/fiber) dengan Dockerfile:

   * `credit-service` — check/reduce/restore limit, webhook consumer for payment
   * `billing-service` — invoice generator, reminder, payment callback handler
   * `ppob-core` — endpoint transaksi (inquiry/pay), uses credit-service & provider connectors
3. Contoh satu `provider-connector` (mock provider + a real-like HTTP adapter) di `ppob-core/pkg/connector`.
4. Database schema (Postgres) SQL migrations for all tables.
5. OpenAPI (YAML) for all public endpoints and sandbox examples.
6. `docker-compose.yml` to run local dev (Postgres, Redis, services, minio for receipts).
7. Unit tests for core logic (≥70% coverage) + integration tests for basic flows.
8. README dengan instruksi build/run/test and Postman collection.

---

## 2. Tech stack & non-fungsional requirement

* **Language:** Golang 1.20+
* **Web framework:** Fiber v2 (or Gin if agent prefers — pilih Fiber untuk konsistensi dengan user request sebelumnya)
* **DB:** PostgreSQL
* **Cache/short-state:** Redis
* **Message queue:** RabbitMQ (dev: use in-memory alternative for tests)
* **Storage:** MinIO (for receipts)
* **Container:** Docker
* **CI:** GitHub Actions (build, lint, test, docker build)
* **Logging:** structured logs (zap or zerolog)
* **Config:** Viper or envconfig (env-driven)
* **Auth:** API Key for tenants (simple header `X-API-Key`), JWT for dashboard (out-of-scope)
* **Idempotency:** `Idempotency-Key` header for transaction create

**SLO (MVP):** handle 500TPS per ppob-core instance; P95 latency < 300ms (local sim ok).

---

## 3. Database schema (migrations) — sketsa

Buat migration SQL files (`migrations/001_init.sql`) dengan tabel utama:

* `tenants` (id, name, plan, api_key, created_at)
* `partners` (id, tenant_id, name, kyc_status, contact, created_at)
* `credit_limits` (partner_id, limit_total, limit_used, limit_available, reset_period, updated_at)
* `transactions` (id, tenant_id, partner_id, amount, fee, total, status, provider, provider_tx_id, idempotency_key, created_at, updated_at)
* `invoices` (id, tenant_id, partner_id, period_start, period_end, amount_due, status, due_date, created_at)
* `receivables` (id, tenant_id, partner_id, amount, status, tx_id, created_at)
* `provider_configs` (id, tenant_id, provider_name, creds_encrypted, endpoint, active)
* `audit_logs` (id, resource_type, resource_id, action, payload_json, created_at)

Add necessary indexes (tenant_id, partner_id, idempotency_key unique on transactions)

---

## 4. API Spec (OpenAPI summary)

Public endpoints (include sample payloads in OpenAPI):

### A. `PPoB Core` (ppob-core)

* `POST /v1/{tenant}/transactions` — create transaction

  * Headers: `X-API-Key`, `Idempotency-Key`
  * Body: `{ "product_code": "PLN_PREPAID", "customer_no": "08123...", "amount": 50000, "partner_id": "p_123" }`
  * Flows:

    1. Validate tenant & api key
    2. Check partner credit via `credit-service` (RPC/HTTP)
    3. Reserve (debit) credit: `reduce_limit(partner_id, amount+fee, tx_id)`
    4. Call provider connector `Pay()` (HTTP)
    5. On success: mark transaction success, create invoice item
    6. On provider failure: refund (restore_limit) and return error
* `GET /v1/{tenant}/transactions/{tx_id}`

### B. `Credit Service` (credit-service)

* `GET /v1/partners/{partner_id}/limit` — returns available limit
* `POST /v1/partners/{partner_id}/reserve` — body `{tx_id, amount}` — atomically reduce available limit; return success/fail
* `POST /v1/partners/{partner_id}/restore` — restore amount (used on refunds)

### C. `Billing Service` (billing-service)

* `GET /v1/tenants/{tenant_id}/invoices` — list
* `POST /v1/tenants/{tenant_id}/invoices/generate` — forced generate
* `POST /v1/payments/callback` — provider/VA payment callback → mark invoice/receivable paid → call credit-service restore or increase

---

## 5. Concurrency, Idempotency, and Atomicity

* `reserve` in credit-service MUST be atomic (DB transaction). If not enough available limit, return 409.
* transactions create path must be idempotent by `idempotency_key`.
* Use DB pessimistic locking or `SELECT ... FOR UPDATE` for `credit_limits` updates.

---

## 6. Provider connector design

* Implement Adapter interface:

```go
type Provider interface {
  Inquiry(ctx context.Context, request InquiryRequest) (InquiryResponse, error)
  Pay(ctx context.Context, request PayRequest) (PayResponse, error)
  Cancel(ctx context.Context, request CancelRequest) (CancelResponse, error)
}
```

* Provide a `mock` provider that simulates success/failure and a `http` adapter example for a fictional provider.
* All connector calls should be cancellable with context timeout and retried with exponential backoff (max 3 attempts) before marking failure.

---

## 7. Error handling & retries

* If provider payment fails temporarily (5xx), retry up to 3 times with backoff. If still fails → restore credit & return error.
* If provider returns business failure (insufficient customer funds / invalid number), restore credit and mark transaction failed.

---

## 8. Tests (must include)

* Unit tests for credit reserve/restore (concurrent reservations to simulate race)
* Integration test for full flow: create transaction -> provider success -> transaction success -> invoice created
* One negative test: create transaction -> provider permanent failure -> ensure credit restored
* Basic contract test for provider connector

---

## 9. Dev & Run (docker-compose)

Compose must include services for local dev:

* postgres:13
* redis
* rabbitmq
* minio
* credit-service
* billing-service
* ppob-core

Provide `Makefile` targets:

* `make build` — build images
* `make up` — docker compose up
* `make test` — run unit + integration tests

---

## 10. CI (GitHub Actions)

Pipeline steps:

* checkout
* setup go
* `go vet` / `golangci-lint run`
* `go test ./...` (with race detector)
* build docker images (optional)

---

## 11. Security & operational notes

* Store provider credentials in environment variables / secrets manager — show example using `.env` for dev
* All internal HTTP calls between services must be authenticated using simple service tokens in headers (e.g. `X-Service-Token`)
* Rate limit per partner: default 3 trx/sec; put limit middleware in ppob-core
* Audit: write audit entry for every reserve/restore/transaction state change

---

## 12. Acceptance criteria (how agent knows it's done)

1. `docker compose up` boots all services and DB migrations applied.
2. Can successfully run the scenario:

   * create tenant & partner (seed script)
   * partner has default daily limit (e.g. 1.000.000)
   * `POST /v1/{tenant}/transactions` with amount 100.000 returns success and decreases available limit
   * after provider success, `GET /v1/{tenant}/transactions/{tx}` shows status `success`
3. Integration tests pass (success + failure path)
4. OpenAPI file exists and endpoints implement the spec
5. README contains run instructions

---

## 13. Coding style & documentation

* Use idiomatic Go. Favor small packages and interfaces.
* Document public functions/types with comments.
* Provide README per service describing env vars and run commands.

---

## 14. Small sample code snippets (starter)

`ppob-core/cmd/main.go` — bootstrap Fiber, config, register routes, connect to credit-service and billing-service via HTTP client (retries & timeouts).

`credit-service/internal/handler/reserve.go` — implement reserve logic (`SELECT FOR UPDATE`, check `limit_available >= amount`, update `limit_used`, `limit_available`).

`billing-service/worker/invoice_generator.go` — cron job that runs daily to generate invoices for active partners.

---

## 15. Notes untuk Agent Copilot

* Prioritaskan correctness & safety over premature optimization.
* Implement minimal mocks to simulate provider behavior for dev & tests.
* Keep secrets out of repo; provide `.env.example` with placeholders.
* Tag each commit clearly (e.g., `feat: credit-service initial`, `test: add reserve tests`).

---

Jika semua sudah jelas, mulai generate repo skeleton, migrations, dan implementasi minimal untuk **reserve** flow terlebih dahulu (credit-service + ppob-core transaction create + mock provider). Setelah itu buat billing worker dan invoice flow.

---

*End of instructions.*
