// @title RediLink URL Shortener API
// @version 1.0
// @description A URL shortening service built with Go and Redis. This API provides endpoints to shorten long URLs and resolve shortened URLs back to their original destinations.
// @description
// @description Features:
// @description - URL shortening with optional custom short codes
// @description - Automatic URL expiration (default 24 hours)
// @description - Rate limiting (10 requests per 30 minutes per IP)
// @description - Redis-based storage for high performance
// @description - Usage analytics tracking
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /
// @schemes http https
// @produces application/json
// @consumes application/json
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/kunalsinghdadhwal/redilink/routes"
	"github.com/watchakorn-18k/scalar-go"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// RateLimitErrorResponse represents a rate limit error response
type RateLimitErrorResponse struct {
	Error          string `json:"error" example:"Rate limit exceeded"`
	RateLimitReset int    `json:"rate_limit_reset" example:"29"`
}

func setupRoutes(app *fiber.App) {
	// API documentation route (must be before catch-all route)
	app.Get("/api/reference", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL:  "./docs/swagger.yaml",
			Theme:    scalar.ThemeId("saturn"),
			DarkMode: true,
		})

		if err != nil {
			return err
		}
		c.Type("html")
		return c.SendString(htmlContent)
	})

	// API routes
	app.Post("/api/v1", routes.ShortenURL)

	// Static file serving for swagger docs
	app.Static("/docs", "./docs")

	// Catch-all route for URL resolution (must be last)
	app.Get("/:url", routes.ResolveURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal((app.Listen(os.Getenv("APP_PORT"))))
}
