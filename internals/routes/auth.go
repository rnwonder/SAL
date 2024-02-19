package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/rnwonder/SAL/data"
	"github.com/rnwonder/SAL/internals/handlers"
	"github.com/rnwonder/SAL/util"
	"github.com/rnwonder/SAL/validators"
)

func AuthRoute(router fiber.Router) {
	router.Post("/", register)
	router.Post("/login", login)
}

func register(ctx *fiber.Ctx) error {
	registeringUser := new(data.MerchantRegister)
	store := ctx.Locals("session").(*session.Store)
	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}

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
			return ctx.Status(400).JSON(fiber.Map{
				"message": "Email already exists",
			})
		}
		if merchant.SkuId == registeringUser.SkuId {
			return ctx.Status(400).JSON(fiber.Map{
				"message": "SkuId is already taken",
			})
		}
	}

	password := util.HashPassword(registeringUser.Password)

	newUser := data.Merchant{
		Email:    registeringUser.Email,
		Password: password,
		SkuId:    registeringUser.SkuId,
		Name:     registeringUser.Name,
		Id:       uuid.Must(uuid.NewRandom()).String(),
	}

	data.MerchantData = append(data.MerchantData, newUser)

	user, token := handlers.LoginUser(sess, &newUser)

	return ctx.Status(200).JSON(fiber.Map{
		"message":   "User registered successfully",
		"user":      user,
		"token":     token,
		"tokenType": "Bearer",
	})

}

func login(ctx *fiber.Ctx) error {
	loginUser := new(data.MerchantLogin)
	store := ctx.Locals("session").(*session.Store)
	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}
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
				user, token := handlers.LoginUser(sess, &merchant)
				return ctx.Status(200).JSON(fiber.Map{
					"message":   "User logged in successfully",
					"user":      user,
					"token":     token,
					"tokenType": "Bearer",
				})
			}
		}
	}

	return ctx.Status(400).JSON(fiber.Map{
		"message": "Invalid credentials",
	})
}
