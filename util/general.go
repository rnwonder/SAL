package util

import (
	"cmp"
	"fmt"
	"github.com/google/uuid"
	"github.com/rnwonder/SAL/internals/models"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

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

func GenerateRandomProducts(numProducts int, merchantId string) map[string]models.Product {
	products := make(map[string]models.Product)
	nameSet := make(map[string]bool)

	for i := 0; i < numProducts; i++ {
		name := generateUniqueName(nameSet)
		id := uuid.Must(uuid.NewRandom()).String()
		products[id] = models.Product{
			SkuId:       merchantId,
			Id:          id,
			Name:        name,
			Description: "Description",
			Price:       rand.Float32() * 100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
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

func SeedData() {
	//models.ProductData = GenerateRandomProducts(40, "merchant1")
}

func EncodeMapToString(data map[string]interface{}) string {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, fmt.Sprintf("%v", value))
	}
	return values.Encode()
}
