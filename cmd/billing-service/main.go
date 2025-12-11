package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aziz46/core-e-voucher-services/internal/billing/handler"
	"github.com/aziz46/core-e-voucher-services/pkg/config"
	"github.com/aziz46/core-e-voucher-services/pkg/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database
	ctx := context.Background()
	err := db.InitDB(ctx, cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, "e_voucher")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Create Fiber app
	app := fiber.New()

	// Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	v1 := app.Group("/v1")

	// Billing endpoints
	v1.Get("/tenants/:tenant_id/invoices", handler.ListInvoices)
	v1.Post("/tenants/:tenant_id/invoices/generate", handler.GenerateInvoice)
	v1.Post("/payments/callback", handler.PaymentCallback)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Billing Service listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
