package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rnwonder/SAL/data"
	"github.com/rnwonder/SAL/util"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

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
				`"totalProducts":32`,
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
				`"totalProducts":32`,
				`"totalPages":7`,
				`"currentPage":2`,
				`"nextPage":"/products?page=3"`,
				`"prevPage":"/products?page=1"`,
			},
			seed: true,
		},
	}

	app := fiber.New()
	ProductRoute(app.Group("/products"))

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
	ProductRoute(app.Group("/products"))

	data.ProductData = []data.Product{
		{
			SkuId:       "37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "A product",
			Description: "A product description",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
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
		isAuth       bool
	}{
		{
			description:  "Create a product with no body",
			route:        "/products",
			expectedCode: 400,
			contains: []string{
				`"message":"Invalid request"`,
			},
			noBody: true,
			isAuth: true,
		},
		{
			description:  "Create a product with invalid body",
			route:        "/products",
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
			isAuth: true,
		},
		{
			description:  "Create a product with valid body",
			route:        "/products",
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
				"skuId":       "shaggsas",
			},
			isAuth: true,
		},
		{
			description:  "Create a product with valid body but no auth",
			route:        "/products",
			expectedCode: 401,
			contains: []string{
				`"message":"Unauthorized"`,
			},
			body: map[string]interface{}{
				"Name":        "A product2",
				"Description": "A product description",
				"Price":       100.00,
			},
		},
		{
			description:  "Create a product with a name that already exists for the merchant",
			route:        "/products",
			expectedCode: 409,
			contains: []string{
				`"message":"You already have a product with this name"`,
			},
			body: map[string]interface{}{
				"Name":        "A product",
				"Description": "A product description",
				"Price":       100.00,
				"skuId":       "ghghjgy",
			},
			isAuth: true,
		},
		{
			description:  "Create a product with a name that already exists for another merchant",
			route:        "/products",
			expectedCode: 201,
			contains: []string{
				`"message":"Product created successfully"`,
			},
			body: map[string]interface{}{
				"Name":        "Car",
				"Description": "A product description",
				"Price":       100.00,
				"skuId":       "ytfgfg",
			},
			isAuth: true,
		},
		{
			description:  "Create a product with a skuId that already exists",
			route:        "/products",
			expectedCode: 409,
			contains: []string{
				`"message":"Product with this SKU already exists"`,
			},
			body: map[string]interface{}{
				"Name":        "Mower",
				"Description": "A product description",
				"Price":       100.00,
				"skuId":       "ytfgfg",
			},
			isAuth: true,
		},
	}

	app := fiber.New()
	AuthRoute(app.Group("/auth"))
	ProductRoute(app.Group("/products"))

	password := util.HashPassword("password")
	data.MerchantData = []data.Merchant{
		{
			Id:        "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:      "Test Merchant",
			Email:     "testMerchant@example.com",
			Password:  password,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	data.ProductData = []data.Product{
		{
			SkuId:       "37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "A product",
			Description: "A product description",
			MerchantId:  "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SkuId:       "56fe9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "Car",
			Description: "A product description",
			MerchantId:  "12rg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, test := range tests {
		body := strings.NewReader(util.EncodeMapToString(test.body))
		var token string
		var req = new(http.Request)

		if test.isAuth {
			authBody := map[string]interface{}{
				"email":    "testMerchant@example.com",
				"password": "password",
			}
			request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(util.EncodeMapToString(authBody)))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(request, 1000)
			if err != nil {
				t.Errorf("error testing route /auth/login: %v", err)
				continue
			}
			read, _ := io.ReadAll(resp.Body)
			jsonString := string(read)
			jsonData := util.JsonParse(jsonString)
			token = jsonData["token"].(string)
		}

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
		isAuth       bool
	}{
		{
			description:  "Delete a product with no auth",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 401,
			contains: []string{
				`"message":"Unauthorized"`,
			},
		},
		{
			description:  "Delete a product that does not exist",
			route:        "/products/dada",
			expectedCode: 404,
			contains: []string{
				`"message":"Product not found"`,
			},
			isAuth: true,
		},
		{
			description:  "Delete a product that exists and is owned by the merchant",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 200,
			contains: []string{
				`"message":"Product deleted successfully"`,
			},
			isAuth: true,
		},
		{
			description:  "Delete a product that exists but is not owned by the merchant",
			route:        "/products/56fe9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 403,
			contains: []string{
				`"message":"You are not authorized to delete this product"`,
			},
			isAuth: true,
		},
	}

	app := fiber.New()
	AuthRoute(app.Group("/auth"))
	ProductRoute(app.Group("/products"))

	password := util.HashPassword("password")
	data.MerchantData = []data.Merchant{
		{
			Id:        "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:      "Test Merchant",
			Email:     "testMerchant@example.com",
			Password:  password,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	data.ProductData = []data.Product{
		{
			SkuId:       "37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "A product",
			Description: "A product description",
			MerchantId:  "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SkuId:       "56fe9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "Car",
			Description: "A product description",
			MerchantId:  "12rg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, test := range tests {
		var token string
		var req = new(http.Request)

		if test.isAuth {
			authBody := map[string]interface{}{
				"email":    "testMerchant@example.com",
				"password": "password",
			}
			request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(util.EncodeMapToString(authBody)))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(request, 1000)
			if err != nil {
				t.Errorf("error testing route /auth/login: %v", err)
				continue
			}
			read, _ := io.ReadAll(resp.Body)
			jsonString := string(read)
			jsonData := util.JsonParse(jsonString)
			token = jsonData["token"].(string)
		}

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
		isAuth       bool
	}{
		{
			description:  "Update a product with no body",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 400,
			contains: []string{
				`"message":"Invalid request"`,
			},
			noBody: true,
			isAuth: true,
		},
		{
			description:  "Update a product with invalid body",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 200,
			contains: []string{
				`"message":"Product updated successfully"`,
			},
			body:   map[string]interface{}{},
			isAuth: true,
		},
		{
			description:  "Update a product with valid body",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 200,
			contains: []string{
				`"message":"Product updated successfully"`,
				`"name":"A product2"`,
			},
			body: map[string]interface{}{
				"Name": "A product2",
			},
			isAuth: true,
		},
		{
			description:  "Update a product with a name that already exists for another product own by the merchant",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 409,
			contains: []string{
				`"message":"You already have a product with this name"`,
			},
			body: map[string]interface{}{
				"Name": "Door",
			},
			isAuth: true,
		},
		{
			description:  "Update a product with a name that already exists for another product not own by the merchant",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 200,
			contains: []string{
				`"message":"Product updated successfully"`,
			},
			body: map[string]interface{}{
				"Name": "Car",
			},
			isAuth: true,
		},
		{
			description:  "Update a product with a skuId that already exists",
			route:        "/products/37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			expectedCode: 409,
			contains: []string{
				`"message":"Product with this SKU already exists"`,
			},
			body: map[string]interface{}{
				"skuId": "56fe9f5c-daf5-46bb-b729-4ceb710f794d",
			},
			isAuth: true,
		},
	}

	app := fiber.New()
	AuthRoute(app.Group("/auth"))
	ProductRoute(app.Group("/products"))

	password := util.HashPassword("password")
	data.MerchantData = []data.Merchant{
		{
			Id:        "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:      "Test Merchant",
			Email:     "testMerchant@example.com",
			Password:  password,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	data.ProductData = []data.Product{
		{
			SkuId:       "37ce9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "A product",
			Description: "A product description",
			MerchantId:  "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SkuId:       "56fe9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "Car",
			Description: "A product description",
			MerchantId:  "12rg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SkuId:       "98je9f5c-daf5-46bb-b729-4ceb710f794d",
			Name:        "Door",
			Description: "A product description",
			MerchantId:  "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			Price:       100.00,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, test := range tests {
		body := strings.NewReader(util.EncodeMapToString(test.body))
		var token string
		var req = new(http.Request)

		if test.isAuth {
			authBody := map[string]interface{}{
				"email":    "testMerchant@example.com",
				"password": "password",
			}
			request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(util.EncodeMapToString(authBody)))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(request, 1000)
			if err != nil {
				t.Errorf("error testing route /auth/login: %v", err)
				continue
			}
			read, _ := io.ReadAll(resp.Body)
			jsonString := string(read)
			jsonData := util.JsonParse(jsonString)
			token = jsonData["token"].(string)
		}

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
