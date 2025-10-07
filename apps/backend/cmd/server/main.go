package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Load .env if present (ignored if missing)
	_ = godotenv.Load()

	// App port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Build DB connection string
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)

	// Connect to Postgres
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to open DB connection: %v", err)
	}
	defer db.Close()

	// Verify connection at startup
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Database not reachable: %v", err)
	}
	log.Println("✅ Connected to Postgres successfully")

	// Initialize Fiber app
	app := fiber.New()

	// --- Enable CORS ---
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // frontend origin
		AllowMethods: "GET,POST,HEAD,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// --- Health check route ---
	app.Get("/healthz", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			log.Println("❌ DB health check failed:", err)
			return c.Status(500).JSON(fiber.Map{"status": "db unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// --- Graceful shutdown ---
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	log.Printf("🚀 Backend running on port %s", port)

	// Wait for Ctrl+C / SIGTERM to stop gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("🛑 Shutting down backend...")
	_ = app.Shutdown()
	log.Println("✅ Server stopped cleanly")
}
