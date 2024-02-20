package main

import (
	"cmp"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/rnwonder/SAL/internals/middleware"
	"github.com/rnwonder/SAL/internals/routes"
	"github.com/rnwonder/SAL/util"
	"os"
)

func main() {
	loadEnvFileError := godotenv.Load("../../.env")

	if loadEnvFileError != nil {
		log.Error("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(cors.New())
	app.Use(middleware.LogRequest)

	// routes
	routes.AuthRoute(app.Group("/auth"))
	routes.ProductRoute(app.Group("/product"))

	app.Get("/", welcomeToApi)
	app.Use(notFound)

	// seed data
	util.SeedData()

	port := cmp.Or(os.Getenv("PORT"), "8000")
	host := cmp.Or(os.Getenv("HOST"), "")

	err := app.Listen(host + ":" + port)
	if err != nil {
		log.Error(err)
	}
}

func welcomeToApi(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "Welcome to the ShopAnythingLagos",
		"live":    true,
	})
}

func notFound(ctx *fiber.Ctx) error {
	return ctx.Status(404).JSON(fiber.Map{
		"message": "Route not found",
	})
}
