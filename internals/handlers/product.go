package handlers

import (
	"cmp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rnwonder/SAL/internals/models"
	"github.com/rnwonder/SAL/types"
	"github.com/rnwonder/SAL/util"
	"github.com/rnwonder/SAL/validators"
	"sync"
	"time"
)

// GetAllProducts Get all products
// @Summary Get all products
// @Description Get all products in the store
// @Tags Product
// @Success 200 {object} GetProductResponse
// @Router /product [get]

func GetAllProductsEndpoint(ctx *fiber.Ctx) error {
	key := cmp.Or(ctx.Query("sortKey"), "createdAt")
	order := cmp.Or(ctx.Query("sortOrder"), "desc")
	searchQuery := ctx.Query("search")

	products := models.GetAllProducts()

	if searchQuery != "" {
		products = models.FilterProductsByName(products, searchQuery)
	}

	resultCh := make(chan models.Product)

	var wg sync.WaitGroup

	numGoroutines := 5

	models.ChunkProductsToChannel(products, resultCh, numGoroutines, &wg)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var resultProducts []models.Product
	for product := range resultCh {
		resultProducts = append(resultProducts, product)
	}

	if key != "" {
		models.SortProducts(resultProducts, key, order)
	}

	startIndex, endIndex, totalPages, limit, page := util.CalculatePageInfo(ctx.Query("page"), ctx.Query("limit"), len(resultProducts))

	resultProducts = resultProducts[startIndex:endIndex]

	return ctx.Status(200).JSON(types.GetProductResponse{
		Message:  "Products fetched successfully",
		Products: resultProducts,
		Meta: types.Meta{
			CurrentPage:   page,
			Limit:         limit,
			TotalPages:    totalPages,
			NextPage:      "/products?page=" + util.NextPage(page, totalPages),
			PrevPage:      "/products?page=" + util.PrevPage(page),
			TotalProducts: len(products),
		},
	})
}

func FindAProductEndpoint(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	product, ok := models.FindProductById(id)

	if !ok {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	return ctx.Status(200).JSON(types.OneProductResponse{
		Message: "Product fetched successfully",
		Product: product,
	})
}

func CreateProductEndpoint(ctx *fiber.Ctx) error {
	body := new(types.ProductCreatePayload)
	skuId := ctx.Query("skuId")

	if skuId == "" {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Invalid request please provide skuId query parameter",
		})
	}

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	if err := validators.Validator(body); err != nil {
		return ctx.Status(400).JSON(err)
	}

	newProduct := models.Product{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		SkuId:       skuId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Id:          uuid.Must(uuid.NewRandom()).String(),
	}

	models.ProductData[newProduct.Id] = newProduct

	return ctx.Status(201).JSON(types.OneProductResponse{
		Message: "Product created successfully",
		Product: newProduct,
	})
}

func UpdateProductEndpoint(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	body := new(types.ProductUpdatePayload)
	skuId := ctx.Query("skuId")

	if skuId == "" {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Invalid request please provide skuId query parameter",
		})
	}

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	if err := validators.Validator(body); err != nil {
		return ctx.Status(400).JSON(err)
	}

	product, ok := models.FindProductById(id)

	if !ok {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	if product.SkuId != skuId {
		return ctx.Status(409).JSON(fiber.Map{
			"message": "You do not have permission to update this product",
		})
	}

	if body.Name != "" {
		product.Name = body.Name
	}

	if body.Description != "" {
		product.Description = body.Description
	}

	if body.Price != 0 {
		product.Price = body.Price
	}

	models.ProductData[product.SkuId] = product

	return ctx.Status(200).JSON(types.OneProductResponse{
		Message: "Product updated successfully",
		Product: product,
	})
}

func DeleteProductEndpoint(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	skuId := ctx.Query("skuId")

	if skuId == "" {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Invalid request please provide skuId query parameter",
		})
	}

	product, ok := models.FindProductById(id)

	if !ok {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	if product.SkuId != skuId {
		return ctx.Status(403).JSON(fiber.Map{
			"message": "You do not have permission to update this product",
		})
	}

	delete(models.ProductData, product.Id)

	return ctx.Status(200).JSON(types.MessageResponse{
		Message: "Product deleted successfully",
	})
}
