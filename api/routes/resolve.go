package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/kunalsinghdadhwal/redilink/database"
)

// ResolveURL godoc
// @Summary Resolve a shortened URL
// @Description Redirect to the original URL using the shortened code. This endpoint also increments the usage counter.
// @Tags URLs
// @Param url path string true "Shortened URL code" minlength(1) maxlength(50) example("abc123")
// @Success 301 "Permanent redirect to the original URL"
// @Failure 404 {object} map[string]string "Not Found - URL not found or expired"
// @Failure 500 {object} map[string]string "Internal Server Error - Database connection issues"
// @Router /{url} [get]
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, fiber.StatusMovedPermanently)
}
