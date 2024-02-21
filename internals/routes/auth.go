package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rnwonder/SAL/data"
	"github.com/rnwonder/SAL/internals/handlers"
	"github.com/rnwonder/SAL/util"
	"github.com/rnwonder/SAL/validators"
)

func AuthRoute(router fiber.Router) {
	router.Post("/register", register)
	router.Post("/login", login)
}

func register(ctx *fiber.Ctx) error {
	registeringUser := new(data.MerchantRegister)

	if err := ctx.BodyParser(registeringUser); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if err := validators.Validator(registeringUser); err != nil {
		return ctx.Status(400).JSON(err)
	}

	for _, merchant := range data.MerchantData {
		if merchant.Email == registeringUser.Email {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "Email already exists",
			})
		}
	}

	password := util.HashPassword(registeringUser.Password)

	newUser := data.Merchant{
		Email:    registeringUser.Email,
		Password: password,
		Name:     registeringUser.Name,
		Id:       uuid.Must(uuid.NewRandom()).String(),
	}

	data.MerchantData = append(data.MerchantData, newUser)

	user, token, expiresAt := handlers.LoginUser(&newUser)

	return ctx.Status(201).JSON(fiber.Map{
		"message":   "User registered successfully",
		"user":      user,
		"token":     token,
		"tokenType": "Bearer",
		"expiresAt": expiresAt,
	})

}

func login(ctx *fiber.Ctx) error {
	loginUser := new(data.MerchantLogin)
	if err := ctx.BodyParser(loginUser); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if err := validators.Validator(loginUser); err != nil {
		return ctx.Status(400).JSON(err)
	}

	for _, merchant := range data.MerchantData {
		if merchant.Email == loginUser.Email {
			if util.CompareHashAndPassword(merchant.Password, loginUser.Password) {
				user, token, expiresAt := handlers.LoginUser(&merchant)
				return ctx.Status(200).JSON(fiber.Map{
					"message":   "User logged in successfully",
					"user":      user,
					"token":     token,
					"tokenType": "Bearer",
					"expiresAt": expiresAt,
				})
			}
		}
	}

	return ctx.Status(401).JSON(fiber.Map{
		"message": "Invalid credentials",
	})
}
