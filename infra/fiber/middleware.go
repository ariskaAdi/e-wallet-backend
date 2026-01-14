package infrafiber

import (
	"ariskaAdi/e-wallet/infra/response"
	"ariskaAdi/e-wallet/internal/config"
	"ariskaAdi/e-wallet/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CheckAuth() fiber.Handler {
	return func (c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		if authorization == "" {
			return NewResponse(
				WithError(response.ErrorUnauthorized),
			).Send(c)
		}

		bearer := strings.Split(authorization, "Bearer ")
		if len(bearer) != 2 {
			log.Println("token invalid")
			return NewResponse(
				WithError(response.ErrorUnauthorized),
			).Send(c)
		}

		token := bearer[1]

		public_id, err := utils.ValidateToken(token, config.Cfg.App.Encryption.JWTSecret)
		if err != nil {
			log.Println("token invalid")
			return NewResponse(
				WithError(response.ErrorUnauthorized),
			).Send(c)
		}

		c.Locals("user_public_id", public_id)
		return c.Next()
	}
}