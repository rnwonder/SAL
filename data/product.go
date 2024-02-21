package data

import "time"

type Product struct {
	SkuId       string    `json:"skuId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	MerchantId  string    `json:"merchantId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProductCreatePayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
	SkuId       string  `json:"skuId" validate:"required"`
}

type ProductUpdatePayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SkuId       string  `json:"skuId"`
}

var ProductData = []Product{}
