package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aziz46/core-e-voucher-services/internal/credit/handler"
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

	// Credit endpoints
	v1.Get("/partners/:partner_id/limit", handler.GetLimit)
	v1.Post("/partners/:partner_id/reserve", handler.Reserve)
	v1.Post("/partners/:partner_id/restore", handler.Restore)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Credit Service listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
