package data

import "time"

type Product struct {
	SkuId       string
	Name        string
	Description string
	Price       float32
	MerchantId  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductResponse struct {
	SkuId       string
	Name        string
	Description string
	Price       float32
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
