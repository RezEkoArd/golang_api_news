package middleware

import (
	"bwanews/config"
	"bwanews/internal/adapter/handler/response"
	"bwanews/lib/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)



type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

// CheckToken implements Middleware
func (o Options) CheckToken() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		authHandler := c.Get("Authorization")

		if authHandler == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Missing Authorization header"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		// fmt.Println("Authorization Header:", authHandler) // Cek header yang diterima

		// Ambil Token dari header "Bearer Token")


		tokenString := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(tokenString)

		if (err) != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid Token"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
			
		}

		c.Locals("user", claims)

		return c.Next()

	}
}


func NewMiddleware(cfg *config.Config) Middleware{
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)

	return opt
}