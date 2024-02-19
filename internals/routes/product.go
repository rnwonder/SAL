package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rnwonder/SAL/data"
	"github.com/rnwonder/SAL/internals/middleware"
	"github.com/rnwonder/SAL/util"
	"github.com/rnwonder/SAL/validators"
	"sort"
	"strings"
	"time"
)

func ProductRoute(router fiber.Router) {
	router.Get("/", getAllProducts)
	router.Post("/", middleware.Authenticated, createAProduct)
	router.Get("/:id", findAProduct)
	router.Put("/:id", middleware.Authenticated, updateAProduct)
	router.Delete("/:id", middleware.Authenticated, deleteAProduct)
}

func getAllProducts(ctx *fiber.Ctx) error {
	var products []data.ProductResponse
	key := util.MyCmpWorkAround(ctx.Query("sortKey"), "createdAt")
	order := util.MyCmpWorkAround(ctx.Query("sortOrder"), "desc")
	searchQuery := ctx.Query("search")

	// Sorting
	sort.Slice(data.ProductData, func(i, j int) bool {
		switch key {
		case "name":
			if order == "asc" {
				return data.ProductData[i].Name < data.ProductData[j].Name
			} else {
				return data.ProductData[i].Name > data.ProductData[j].Name
			}
		case "price":
			if order == "asc" {
				return data.ProductData[i].Price < data.ProductData[j].Price
			} else {
				return data.ProductData[i].Price > data.ProductData[j].Price
			}
		case "createdAt":
			if order == "asc" {
				return data.ProductData[i].CreatedAt.Before(data.ProductData[j].CreatedAt)
			} else {
				return data.ProductData[i].CreatedAt.After(data.ProductData[j].CreatedAt)
			}
		default:
			return i < j
		}
	})

	// Filtering by search query
	filteredData := data.ProductData
	if searchQuery != "" {
		filteredData = []data.Product{}
		for _, product := range data.ProductData {
			if strings.Contains(strings.ToLower(product.Name), strings.ToLower(searchQuery)) {
				filteredData = append(filteredData, product)
			}
		}
	}

	startIndex, endIndex, totalPages, limit, page := util.CalculatePageInfo(ctx.Query("page"), ctx.Query("limit"), len(filteredData))

	for i, product := range filteredData[startIndex:endIndex] {
		products = append(products, util.ClientProductFormat(product))
		if i == endIndex {
			break
		}
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message":  "Products fetched successfully",
		"products": products,
		"meta": fiber.Map{
			"currentPage":   page,
			"limit":         limit,
			"totalPages":    totalPages,
			"nextPage":      "/products?page=" + util.NextPage(page, totalPages),
			"prevPage":      "/products?page=" + util.PrevPage(page),
			"totalProducts": len(filteredData),
		},
	})
}

func createAProduct(ctx *fiber.Ctx) error {
	product := new(data.ProductCreatePayload)

	user := ctx.Locals("user").(data.Merchant)

	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if err := validators.Validator(product); err != nil {
		return ctx.Status(400).JSON(err)
	}

	for _, savedProduct := range data.ProductData {
		if savedProduct.Name == product.Name && savedProduct.SkuId == user.SkuId {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "You already have a product with this name",
			})
		}
	}

	newProduct := data.Product{
		Id:          uuid.Must(uuid.NewRandom()).String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SkuId:       user.SkuId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	data.ProductData = append(data.ProductData, newProduct)

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": util.ClientProductFormat(newProduct),
	})
}

func findAProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	for _, product := range data.ProductData {
		if product.Id == id {
			return ctx.Status(200).JSON(fiber.Map{
				"message": "Product fetched successfully",
				"product": util.ClientProductFormat(product),
			})
		}
	}

	return ctx.Status(404).JSON(fiber.Map{
		"message": "Product not found",
	})
}

func updateAProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	product := new(data.ProductUpdatePayload)
	user := ctx.Locals("user").(data.Merchant)

	if err := ctx.BodyParser(product); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if err := validators.Validator(product); err != nil {
		return ctx.Status(400).JSON(err)
	}

	for i, savedProduct := range data.ProductData {
		if savedProduct.Id == id && savedProduct.SkuId == user.SkuId {

			if product.Name != "" {
				data.ProductData[i].Name = product.Name
			}

			if product.Description != "" {
				data.ProductData[i].Description = product.Description
			}

			if product.Price != 0 {
				data.ProductData[i].Price = product.Price
			}

			data.ProductData[i].UpdatedAt = time.Now()

			return ctx.Status(200).JSON(fiber.Map{
				"message": "Product updated successfully",
				"product": util.ClientProductFormat(data.ProductData[i]),
			})
		}
	}
	return ctx.Status(404).JSON(fiber.Map{
		"message": "Product not found",
	})
}

func deleteAProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user := ctx.Locals("user").(data.Merchant)

	for i, product := range data.ProductData {
		if user.SkuId == product.SkuId && product.Id == id {
			data.ProductData = append(data.ProductData[:i], data.ProductData[i+1:]...)
			return ctx.Status(200).JSON(fiber.Map{
				"message": "Product deleted successfully",
			})
		}
	}
	return ctx.Status(404).JSON(fiber.Map{
		"message": "Product not found",
	})
}
