package main

import (
	"cmp"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	_ "github.com/rnwonder/SAL/docs"
	"github.com/rnwonder/SAL/internals/handlers"
	"github.com/rnwonder/SAL/internals/middleware"
	"github.com/rnwonder/SAL/util"
	"os"
)

// @title           ShopAnythingLagos API
// @version         1.0
// @description     This is the ShopAnythingLagos API documentation
// @host      localhost:4500
// @BasePath  /

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

	products := app.Group("/product")
	products.Get("/", handlers.GetAllProductsEndpoint)
	products.Get("/:id", handlers.FindAProductEndpoint)
	products.Post("/", handlers.CreateProductEndpoint)
	products.Put("/:id", handlers.UpdateProductEndpoint)
	products.Delete("/:id", handlers.DeleteProductEndpoint)

	app.Get("/swagger/*", swagger.HandlerDefault)

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
