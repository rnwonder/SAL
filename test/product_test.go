package test

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rnwonder/SAL/internals/handlers"
	"github.com/rnwonder/SAL/internals/models"
	"github.com/rnwonder/SAL/util"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var testId1 = "37ce9f5c-daf5-46bb-b729-4ceb710f794d"
var testId2 = "56fe9f5c-daf5-46bb-b729-4ceb710f794d"
var testId3 = "98je9f5c-daf5-46bb-b729-4ceb710f794d"

func Test_getAllProducts(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
		seed         bool
	}{
		{
			description:  "Get all products with no seed data",
			route:        "/products",
			expectedCode: 200,
			contains: []string{
				`"products":null`,
				`"message":"Products fetched successfully"`,
				`"limit":10`,
			},
		},
		{
			description:  "Limit the number of products to 5 with no seed data",
			route:        "/products?limit=5",
			expectedCode: 200,
			contains: []string{
				`"products":null`,
				`"message":"Products fetched successfully"`,
				`"limit":5`,
			},
		},
		{
			description:  "Get all products with seed data",
			route:        "/products",
			expectedCode: 200,
			contains: []string{
				`"products":[`,
				`"message":"Products fetched successfully"`,
				`"limit":10`,
			},
			seed: true,
		},
		{
			description:  "Get all products with seed data",
			route:        "/products?limit=5&page=2",
			expectedCode: 200,
			contains: []string{
				`"products":[`,
				`"message":"Products fetched successfully"`,
				`"limit":5`,
				`"currentPage":2`,
				`"nextPage":"/products?page=3"`,
				`"prevPage":"/products?page=1"`,
			},
			seed: true,
		},
	}

	app := fiber.New()
	products := app.Group("/products")
	products.Get("/", handlers.GetAllProductsEndpoint)

	hasSeed := false

	for _, test := range tests {

		if test.seed && !hasSeed {
			util.SeedData()
			hasSeed = true
		}

		req := httptest.NewRequest("GET", test.route, nil)

		resp, err := app.Test(req, 1000)

		if err != nil {
			t.Errorf("error testing route %s: %v", test.route, err)
			continue
		}

		read, _ := io.ReadAll(resp.Body)

		if len(test.contains) > 0 {
			for _, contain := range test.contains {
				assert.Containsf(t, string(read), contain, test.description)
			}
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func Test_findAProduct(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
	}{
		{
			description:  "Find a product that does not exist",
			route:        "/products/dada",
			expectedCode: 404,
			contains: []string{
				`"message":"Product not found"`,
			},
		},
		{
			description:  "Find a product that exists",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 200,
			contains: []string{
				`"message":"Product fetched successfully"`,
				`"name":"A product"`,
				`"price":100`,
				`"description":"A product description"`,
			},
		},
	}

	app := fiber.New()
	products := app.Group("/products")
	products.Get("/:id", handlers.FindAProductEndpoint)

	models.ProductData[testId1] = models.Product{
		Id:          testId1,
		SkuId:       "someSkuId",
		Name:        "A product",
		Description: "A product description",
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, test := range tests {

		req := httptest.NewRequest("GET", test.route, nil)

		resp, err := app.Test(req, 1000)

		if err != nil {
			t.Errorf("error testing route %s: %v", test.route, err)
			continue
		}

		read, _ := io.ReadAll(resp.Body)

		if len(test.contains) > 0 {
			for _, contain := range test.contains {
				assert.Containsf(t, string(read), contain, test.description)
			}
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func Test_createAProduct(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
		body         map[string]interface{}
		noBody       bool
	}{
		{
			description:  "Create a product with no body",
			route:        "/products?skuId=shaggsas",
			expectedCode: 400,
			contains: []string{
				`"message":"Invalid request payload"`,
			},
			noBody: true,
		},
		{
			description:  "Create a product with invalid body",
			route:        "/products?skuId=shaggsas",
			expectedCode: 400,
			contains: []string{
				`Name`,
				`Description`,
				`Price`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"sss": "A product2",
			},
		},
		{
			description:  "Create a product with valid body",
			route:        "/products?skuId=shaggsas",
			expectedCode: 201,
			contains: []string{
				`"message":"Product created successfully"`,
				`"name":"A product2"`,
				`"price":100`,
				`"description":"A product description"`,
			},
			body: map[string]interface{}{
				"Name":        "A product2",
				"Description": "A product description",
				"Price":       100.00,
			},
		},
		{
			description:  "Create a product with valid body but no auth",
			route:        "/products",
			expectedCode: 401,
			contains: []string{
				`"message":"Invalid request please provide skuId query parameter"`,
			},
			body: map[string]interface{}{
				"Name":        "A product2",
				"Description": "A product description",
				"Price":       100.00,
			},
		},
	}

	app := fiber.New()
	products := app.Group("/products")
	products.Post("/", handlers.CreateProductEndpoint)

	models.ProductData[testId1] = models.Product{
		SkuId:       "someSkuId",
		Name:        "A product",
		Description: "A product description",
		Id:          testId1,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	models.ProductData[testId2] = models.Product{
		SkuId:       "someSkuId2",
		Name:        "Car",
		Description: "A product description",
		Id:          testId2,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, test := range tests {
		body := strings.NewReader(util.EncodeMapToString(test.body))
		var token string
		var req = new(http.Request)

		if !test.noBody {
			req = httptest.NewRequest("POST", test.route, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req = httptest.NewRequest("POST", test.route, nil)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, 1000)

		if err != nil {
			t.Errorf("error testing route %s: %v", test.route, err)
			continue
		}

		read, _ := io.ReadAll(resp.Body)

		if len(test.contains) > 0 {
			for _, contain := range test.contains {
				assert.Containsf(t, string(read), contain, test.description)
			}
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func Test_deleteAProduct(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
	}{
		{
			description:  "Delete a product with no auth",
			route:        "/products/" + testId1,
			expectedCode: 401,
			contains: []string{
				`"message":"Invalid request please provide skuId query parameter"`,
			},
		},
		{
			description:  "Delete a product that does not exist",
			route:        "/products/dada?skuId=someSkuId",
			expectedCode: 404,
			contains: []string{
				`"message":"Product not found"`,
			},
		},
		{
			description:  "Delete a product that exists and is owned by the merchant",
			route:        "/products/" + testId1 + "?skuId=someSkuId",
			expectedCode: 200,
			contains: []string{
				`"message":"Product deleted successfully"`,
			},
		},
		{
			description:  "Delete a product that exists but is not owned by the merchant",
			route:        "/products/" + testId2 + "?skuId=someSkuId",
			expectedCode: 403,
			contains: []string{
				`"message":"You do not have permission to update this product"`,
			},
		},
	}

	app := fiber.New()
	products := app.Group("/products")
	products.Delete("/:id", handlers.DeleteProductEndpoint)

	models.ProductData[testId1] = models.Product{
		SkuId:       "someSkuId",
		Name:        "A product",
		Description: "A product description",
		Id:          testId1,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	models.ProductData[testId2] = models.Product{
		SkuId:       "someSkuId2",
		Name:        "Car",
		Description: "A product description",
		Id:          testId2,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, test := range tests {
		var token string
		var req = new(http.Request)

		req = httptest.NewRequest("DELETE", test.route, nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, 1000)

		if err != nil {
			t.Errorf("error testing route %s: %v", test.route, err)
			continue
		}

		read, _ := io.ReadAll(resp.Body)

		if len(test.contains) > 0 {
			for _, contain := range test.contains {
				assert.Containsf(t, string(read), contain, test.description)
			}
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}

func Test_updateAProduct(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
		body         map[string]interface{}
		noBody       bool
	}{
		{
			description:  "Update a product with no body",
			route:        "/products/" + testId1 + "?skuId=someSkuId",
			expectedCode: 400,
			contains: []string{
				`"message":"Invalid request payload"`,
			},
			noBody: true,
		},
		{
			description:  "Update a product with invalid body",
			route:        "/products/" + testId1 + "?skuId=someSkuId",
			expectedCode: 200,
			contains: []string{
				`"message":"Product updated successfully"`,
			},
			body: map[string]interface{}{},
		},
		{
			description:  "Update a product with valid body",
			route:        "/products/" + testId1 + "?skuId=someSkuId",
			expectedCode: 200,
			contains: []string{
				`"message":"Product updated successfully"`,
				`"name":"A product2"`,
			},
			body: map[string]interface{}{
				"Name": "A product2",
			},
		},
		{
			description:  "Update a product with valid body",
			route:        "/products/" + testId1 + "?skuId=someSkuId",
			expectedCode: 200,
			contains: []string{
				`"message":"Product updated successfully"`,
			},
			body: map[string]interface{}{
				"Name": "Car",
			},
		},
	}

	app := fiber.New()
	products := app.Group("/products")
	products.Put("/:id", handlers.UpdateProductEndpoint)

	models.ProductData[testId1] = models.Product{
		SkuId:       "someSkuId",
		Name:        "A product",
		Description: "A product description",
		Id:          testId1,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	models.ProductData[testId2] = models.Product{
		SkuId:       "someSkuId2",
		Name:        "Car",
		Description: "A product description",
		Id:          testId2,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	models.ProductData[testId1] = models.Product{
		SkuId:       "someSkuId",
		Name:        "Door",
		Description: "A product description",
		Id:          testId1,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, test := range tests {
		body := strings.NewReader(util.EncodeMapToString(test.body))
		var token string
		var req = new(http.Request)

		if !test.noBody {
			req = httptest.NewRequest("PUT", test.route, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req = httptest.NewRequest("PUT", test.route, nil)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := app.Test(req, 1000)

		if err != nil {
			t.Errorf("error testing route %s: %v", test.route, err)
			continue
		}

		read, _ := io.ReadAll(resp.Body)

		if len(test.contains) > 0 {
			for _, contain := range test.contains {
				assert.Containsf(t, string(read), contain, test.description)
			}
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
