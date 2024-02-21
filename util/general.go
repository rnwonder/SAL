package util

import (
	"cmp"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/rnwonder/SAL/data"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

func JsonParse(jsonString string) map[string]interface{} {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Error("Error parsing json", err)
		return make(map[string]interface{})
	}
	return data
}

func HashPassword(password string) string {
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func CompareHashAndPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func CalculatePageInfo(page string, limit string, total int) (int, int, int, int, int) {
	pageString := cmp.Or(page, "1")
	limitString := cmp.Or(limit, "10")
	pageInt, _ := strconv.Atoi(pageString)
	limitInt, _ := strconv.Atoi(limitString)
	startIndex := (pageInt - 1) * limitInt
	endIndex := pageInt * limitInt
	totalPages := total / limitInt

	// Account for remainder
	if total%limitInt > 0 {
		totalPages++
	}

	if totalPages < 1 {
		totalPages = 1
	}

	if endIndex > total {
		endIndex = total
	}
	return startIndex, endIndex, totalPages, limitInt, pageInt
}

func NextPage(page int, totalPages int) string {
	if page >= totalPages {
		return strconv.Itoa(totalPages)
	}
	return strconv.Itoa(page + 1)
}

func PrevPage(page int) string {
	if page <= 1 {
		return "1"
	}
	return strconv.Itoa(page - 1)
}

func GenerateRandomProducts(numProducts int, merchantId string) []data.Product {
	var products []data.Product
	nameSet := make(map[string]bool)

	for i := 0; i < numProducts; i++ {
		name := generateUniqueName(nameSet)
		products = append(products, data.Product{
			SkuId:       uuid.Must(uuid.NewRandom()).String(),
			MerchantId:  merchantId,
			Name:        name,
			Description: "Description",
			Price:       rand.Float32() * 100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return products
}

func generateUniqueName(nameSet map[string]bool) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nameLength := 8
	for {
		var nameBuilder string
		for i := 0; i < nameLength; i++ {
			nameBuilder += string(chars[rand.Intn(len(chars))])
		}
		if _, ok := nameSet[nameBuilder]; !ok {
			nameSet[nameBuilder] = true
			return nameBuilder
		}
	}
}

func GenerateRandomMerchants(numMerchants int) []data.Merchant {
	var merchants []data.Merchant
	emailSet := make(map[string]bool)

	for i := 0; i < numMerchants; i++ {
		email := generateUniqueEmail(emailSet)
		merchants = append(merchants, data.Merchant{
			Id:        uuid.Must(uuid.NewRandom()).String(),
			Name:      fmt.Sprintf("Merchant%d", i+1),
			Email:     email,
			Password:  HashPassword("password"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return merchants
}

func generateUniqueEmail(emailSet map[string]bool) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	emailLength := 8
	for {
		var emailBuilder string
		for i := 0; i < emailLength; i++ {
			emailBuilder += string(chars[rand.Intn(len(chars))])
		}
		emailBuilder += "@example.com"
		if _, ok := emailSet[emailBuilder]; !ok {
			emailSet[emailBuilder] = true
			return emailBuilder
		}
	}
}

func SeedData() {
	data.MerchantData = append(data.MerchantData, GenerateRandomMerchants(8)...)

	for _, merchant := range data.MerchantData {
		data.ProductData = append(data.ProductData, GenerateRandomProducts(4, merchant.Id)...)
	}
}

func EncodeMapToString(data map[string]interface{}) string {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}
