package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aziz46/core-e-voucher-services/internal/ppob/handler"
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

	// PPOB endpoints
	v1.Post("/:tenant/transactions", handler.CreateTransaction)
	v1.Get("/:tenant/transactions/:tx_id", handler.GetTransaction)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("PPOB Core Service listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
