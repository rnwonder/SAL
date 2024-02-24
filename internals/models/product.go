package models

import (
	"sort"
	"strings"
	"sync"
	"time"
)

type Product struct {
	SkuId       string    `json:"skuId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var ProductData = map[string]Product{
	"1": {
		Id:          "1",
		SkuId:       "someSkuId",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
		Name:        "Product 1",
		Description: "Description",
		Price:       50,
	},
	"2": {
		Id:          "2",
		SkuId:       "someSkuId",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
		Name:        "Product 2",
		Description: "Description",
		Price:       150,
	},
}

func GetAllProducts() []Product {
	products := make([]Product, 0, len(ProductData))
	for _, product := range ProductData {
		products = append(products, product)
	}
	return products
}

func FilterProductsByName(products []Product, name string) []Product {
	filtered := make([]Product, 0)
	for _, product := range products {
		if strings.Contains(strings.ToLower(product.Name), strings.ToLower(name)) {
			filtered = append(filtered, product)
		}
	}
	return filtered
}

func SortProducts(products []Product, sortBy string, sortOrder string) {
	switch sortBy {
	case "name":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Name < products[j].Name
			}
			return products[i].Name > products[j].Name
		})
	case "price":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Price < products[j].Price
			}
			return products[i].Price > products[j].Price
		})
	case "createdAt":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].CreatedAt.Before(products[j].CreatedAt)
			}
			return products[i].CreatedAt.After(products[j].CreatedAt)
		})
	}
}

func ChunkProductsToChannel(products []Product, channel chan Product, numberOfGoroutines int, wg *sync.WaitGroup) {
	chunkSize := (len(products) + numberOfGoroutines - 1) / numberOfGoroutines
	for i := 0; i < len(products); i += chunkSize {
		wg.Add(1)
		go addProductToChannel(products, channel, i, i+chunkSize, wg)
	}
}

func addProductToChannel(products []Product, channel chan Product, start int, end int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, product := range products[start:end] {
		channel <- product
	}
}

func FindProductById(id string) (Product, bool) {
	product, ok := ProductData[id]
	return product, ok
}
