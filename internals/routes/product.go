package routes

import (
	"github.com/gofiber/fiber/v2"
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
	var products []data.Product
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
		products = append(products, product)
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
		if savedProduct.SkuId == product.SkuId {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "Product with this SKU already exists",
			})
		}
		if savedProduct.Name == product.Name && savedProduct.MerchantId == user.Id {
			return ctx.Status(409).JSON(fiber.Map{
				"message": "You already have a product with this name",
			})
		}
	}

	newProduct := data.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SkuId:       product.SkuId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		MerchantId:  user.Id,
	}

	data.ProductData = append(data.ProductData, newProduct)

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": newProduct,
	})
}

func findAProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	for _, product := range data.ProductData {
		if product.SkuId == id {
			return ctx.Status(200).JSON(fiber.Map{
				"message": "Product fetched successfully",
				"product": product,
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
		if savedProduct.SkuId == id && savedProduct.MerchantId == user.Id {
			if product.Name != "" {
				for _, savedProduct := range data.ProductData {
					if savedProduct.Name == product.Name && savedProduct.MerchantId == user.Id && savedProduct.SkuId != id {
						return ctx.Status(409).JSON(fiber.Map{
							"message": "You already have a product with this name",
						})
					}
				}
				data.ProductData[i].Name = product.Name
			}

			if product.Description != "" {
				data.ProductData[i].Description = product.Description
			}

			if product.Price != 0 {
				data.ProductData[i].Price = product.Price
			}
			if product.SkuId != "" {

				for _, savedProduct := range data.ProductData {
					if savedProduct.SkuId == product.SkuId && savedProduct.SkuId != id {
						return ctx.Status(409).JSON(fiber.Map{
							"message": "Product with this SKU already exists",
						})
					}
				}
				data.ProductData[i].SkuId = product.SkuId
			}

			data.ProductData[i].UpdatedAt = time.Now()

			return ctx.Status(200).JSON(fiber.Map{
				"message": "Product updated successfully",
				"product": data.ProductData[i],
			})
		}
		if savedProduct.SkuId == id && savedProduct.MerchantId != user.Id {
			return ctx.Status(403).JSON(fiber.Map{
				"message": "You are not authorized to update this product",
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
		if user.Id == product.MerchantId && product.SkuId == id {
			data.ProductData = append(data.ProductData[:i], data.ProductData[i+1:]...)
			return ctx.Status(200).JSON(fiber.Map{
				"message": "Product deleted successfully",
			})
		}
		if user.Id != product.MerchantId && product.SkuId == id {
			return ctx.Status(403).JSON(fiber.Map{
				"message": "You are not authorized to delete this product",
			})
		}
	}
	return ctx.Status(404).JSON(fiber.Map{
		"message": "Product not found",
	})
}
