package routes

import (
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

func Test_login(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
		body         map[string]interface{}
		noBody       bool
	}{
		{
			description:  "Login with correct credentials",
			route:        "/auth/login",
			expectedCode: 200,
			contains: []string{
				`"token"`,
				`"tokenType":"Bearer"`,
				`"name":"Test Merchant"`,
				`"email":"testMerchant@example.com"`,
				`"expiresAt"`,
				`"skuId":"someskuid"`,
			},
			body: map[string]interface{}{
				"email":    "testMerchant@example.com",
				"password": "password",
			},
		},
		{
			description:  "Login with incorrect credentials",
			route:        "/auth/login",
			expectedCode: 401,
			contains: []string{
				`"message":"Invalid credentials"`,
			},
			body: map[string]interface{}{
				"email":    "testMerchant@example.com",
				"password": "asasasas",
			},
		},
		{
			description:  "Login with missing email",
			route:        "/auth/login",
			expectedCode: 400,
			contains: []string{
				`Email`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"password": "password",
			},
		},
		{
			description:  "Login with missing password",
			route:        "/auth/login",
			expectedCode: 400,
			contains: []string{
				`Password`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"email": "testMerchant@example.com",
			},
		},
		{
			description:  "Login with missing email and password",
			route:        "/auth/login",
			expectedCode: 400,
			contains: []string{
				`Email`,
				`Needs to implement 'required'`,
				`Password`,
			},
			body: map[string]interface{}{},
		},
		{
			description:  "Login without body",
			route:        "/auth/login",
			expectedCode: 400,
			contains: []string{
				`"message":"Invalid request"`,
			},
			noBody: true,
		},
	}

	app := fiber.New()
	AuthRoute(app.Group("/auth"))

	password := util.HashPassword("password")
	data.MerchantData = []data.Merchant{
		{
			Id:        "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			SkuId:     "someskuid",
			Name:      "Test Merchant",
			Email:     "testMerchant@example.com",
			Password:  password,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	for _, test := range tests {
		body := strings.NewReader(util.EncodeMapToString(test.body))

		var req = new(http.Request)

		if !test.noBody {
			req = httptest.NewRequest("POST", test.route, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req = httptest.NewRequest("POST", test.route, nil)
		}

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

func Test_register(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
		contains     []string
		body         map[string]interface{}
		noBody       bool
	}{
		{
			description: "Register with correct credentials",
			route:       "/auth/register",
			contains: []string{
				`"token"`,
				`"tokenType":"Bearer"`,
				`"name":"Test Merchant"`,
				`"email":"testMerchant@example.com"`,
				`"expiresAt"`,
				`"skuId":"someskuid"`,
			},
			expectedCode: 201,
			body: map[string]interface{}{
				"email":    "testMerchant@example.com",
				"name":     "Test Merchant",
				"skuId":    "someskuid",
				"password": "123456",
			},
		},
		{
			description:  "Register with missing email",
			route:        "/auth/register",
			expectedCode: 400,
			contains: []string{
				`Email`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"name":     "Test Merchant",
				"skuId":    "someskuid",
				"password": "123456",
			},
		},
		{
			description:  "Register with missing name",
			route:        "/auth/register",
			expectedCode: 400,
			contains: []string{
				`Name`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"email":    "testMerchant@example.com",
				"skuId":    "someskuid",
				"password": "123456",
			},
		},
		{
			description:  "Register with missing skuId",
			route:        "/auth/register",
			expectedCode: 400,
			contains: []string{
				`SkuId`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"email":    "testMerchant@example.com",
				"name":     "Test Merchant",
				"password": "123456",
			},
		},
		{
			description:  "Register with missing password",
			route:        "/auth/register",
			expectedCode: 400,
			contains: []string{
				`Password`,
				`Needs to implement 'required'`,
			},
			body: map[string]interface{}{
				"email": "testMerchant@example.com",
				"name":  "Test Merchant",
				"skuId": "someskuid",
			},
		},
		{
			description:  "Register with missing email, name, skuId and password",
			route:        "/auth/register",
			expectedCode: 400,
			contains: []string{
				`Email`,
				`Needs to implement 'required'`,
				`Name`,
				`SkuId`,
				`Password`,
			},
			body: map[string]interface{}{},
		},
		{
			description:  "Register without body",
			route:        "/auth/register",
			expectedCode: 400,
			contains: []string{
				`"message":"Invalid request"`,
			},
			noBody: true,
		},
		{
			description:  "Register with existing email",
			route:        "/auth/register",
			expectedCode: 409,
			contains: []string{
				`"message":"Email already exists"`,
			},
			body: map[string]interface{}{
				"email":    "testMerchant2@example.com",
				"name":     "Test Merchant",
				"skuId":    "someskuid",
				"password": "123456",
			},
		},
		{
			description:  "Register with existing skuId",
			route:        "/auth/register",
			expectedCode: 409,
			contains: []string{
				`"message":"SkuId is already taken"`,
			},
			body: map[string]interface{}{
				"email":    "testMerchant@example.com",
				"name":     "Test Merchant",
				"skuId":    "someskuid2",
				"password": "123456",
			},
		},
	}

	app := fiber.New()
	AuthRoute(app.Group("/auth"))

	password := util.HashPassword("password")
	data.MerchantData = []data.Merchant{
		{
			Id:        "6grg9f5c-daf5-46bb-b729-4ceb710f794d",
			SkuId:     "someskuid2",
			Name:      "Test Merchant2",
			Email:     "testMerchant2@example.com",
			Password:  password,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
	}

	for _, test := range tests {
		body := strings.NewReader(util.EncodeMapToString(test.body))

		var req = new(http.Request)

		if !test.noBody {
			req = httptest.NewRequest("POST", test.route, body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req = httptest.NewRequest("POST", test.route, nil)
		}

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
