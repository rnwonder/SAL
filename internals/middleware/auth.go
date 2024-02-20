package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/rnwonder/SAL/internals/handlers"
	"strings"
)

func Authenticated(ctx *fiber.Ctx) error {
	var token string
	var apiKey string

	token = ctx.Get("Authorization")

	if token == "" {
		log.Error("No token provided")
		err := ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
			"code":    401,
		})
		if err != nil {
			return err
		}
		return nil
	}

	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		log.Error("Invalid/Malformed auth token")
		err := ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
			"code":    401,
		})
		if err != nil {
			return err
		}
		return nil
	}
	apiKey = parts[1]

	userData, err := handlers.DecodeJWTData(apiKey)
	if err != nil {
		log.Error("Error decoding token", err)
		err := ctx.Status(401).JSON(fiber.Map{
			"message": "Expired or invalid api key",
			"code":    401,
		})
		if err != nil {
			return err
		}
		return nil
	}

	ctx.Locals("user", userData)
	nextError := ctx.Next()

	if nextError != nil {
		log.Error("Error in next", nextError)
		err := ctx.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
			"code":    500,
		})
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
