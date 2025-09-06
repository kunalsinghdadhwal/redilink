package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kunalsinghdadhwal/redilink/database"
	"github.com/kunalsinghdadhwal/redilink/helpers"
)

// ShortenRequest represents the request body for shortening a URL
type request struct {
	URL         string `json:"url" example:"https://www.example.com" validate:"required,url"` // The original URL to be shortened (required)
	CustomShort string `json:"short,omitempty" example:"mylink"`                              // Optional custom short code (if not provided, a random one will be generated)
	Expiry      int    `json:"expiry,omitempty" example:"24"`                                 // URL expiration time in hours (default: 24 hours)
}

// ShortenResponse represents the response body for URL shortening
type response struct {
	URL                string `json:"url" example:"https://www.example.com"`        // The original URL that was shortened
	CustomShort        string `json:"short" example:"http://localhost:3000/abc123"` // The complete shortened URL
	Expiry             int    `json:"expiry" example:"24"`                          // URL expiration time in hours
	X_Rate_Remaining   int    `json:"rate_limit" example:"9"`                       // Number of requests remaining for current IP
	X_Rate_Limit_Reset int    `json:"rate_limit_reset" example:"29"`                // Time until rate limit resets (in minutes)
}

// ShortenURL godoc
// @Summary Shorten a URL
// @Description Create a shortened URL from a given long URL. Rate limited to 10 requests per 30 minutes per IP.
// @Tags URLs
// @Accept json
// @Produce json
// @Param request body request true "URL shortening request"
// @Success 200 {object} response "Successfully shortened URL"
// @Failure 400 {object} map[string]string "Bad Request - Invalid request body or invalid URL"
// @Failure 403 {object} map[string]string "Forbidden - Custom short URL already exists"
// @Failure 429 {object} map[string]interface{} "Too Many Requests - Rate limit exceeded"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Failure 503 {object} map[string]string "Service Unavailable - Domain not allowed"
// @Router /api/v1 [post]
func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": int(limit / time.Nanosecond / time.Minute),
			})
		}
	}
	// Validate URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Speak to me, Oh Toothless One"})
	}

	body.URL = helpers.EnforeHTTP(body.URL)

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()

	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom short URL already exists",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, time.Duration(body.Expiry)*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	resp := response{
		URL:                body.URL,
		CustomShort:        "",
		Expiry:             body.Expiry,
		X_Rate_Remaining:   10,
		X_Rate_Limit_Reset: 30,
	}

	r2.Decr(database.Ctx, c.IP())
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.X_Rate_Remaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.X_Rate_Limit_Reset = int(ttl / time.Nanosecond / time.Minute)
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
