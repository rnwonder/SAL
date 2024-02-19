package util

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/rnwonder/SAL/data"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func JsonStringify(data map[string]interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("Error converting data to json", err)
		return ""
	}
	jsonString := string(jsonData)
	return jsonString
}

func JsonParse(jsonString string) map[string]interface{} {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		log.Error("Error parsing json", err)
		return make(map[string]interface{})
	}
	return data
}

func MyCmpWorkAround(value1 string, value2 string) string {
	// This is a work around for the cmp package breaking in 1.21.6
	if value1 == "" {
		return value2
	}
	return value1
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

func CalculatePageInfo(page string, limit string, total int) (int, int, int) {
	pageString := MyCmpWorkAround(page, "1")
	limitString := MyCmpWorkAround(limit, "10")
	pageInt, _ := strconv.Atoi(pageString)
	limitInt, _ := strconv.Atoi(limitString)
	startIndex := (pageInt - 1) * limitInt
	endIndex := pageInt * limitInt
	totalPages := total / limitInt

	if endIndex > total {
		endIndex = total
	}
	return startIndex, endIndex, totalPages
}

func NextPage(page string, totalPages int) string {
	pageString := MyCmpWorkAround(page, "1")
	pageInt, _ := strconv.Atoi(pageString)
	if pageInt >= totalPages {
		return strconv.Itoa(totalPages)
	}
	return strconv.Itoa(pageInt + 1)
}

func PrevPage(page string) string {
	pageString := MyCmpWorkAround(page, "1")
	pageInt, _ := strconv.Atoi(pageString)
	if pageInt <= 1 {
		return "1"
	}
	return strconv.Itoa(pageInt - 1)
}

func ClientProductFormat(product data.Product) data.ProductResponse {
	return data.ProductResponse{
		Id:          product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
