package types

import "github.com/rnwonder/SAL/internals/models"

type Meta struct {
	CurrentPage   int    `json:"currentPage"`
	Limit         int    `json:"limit"`
	TotalPages    int    `json:"totalPages"`
	NextPage      string `json:"nextPage"`
	PrevPage      string `json:"prevPage"`
	TotalProducts int    `json:"totalProducts"`
}

type GetProductResponse struct {
	Products []models.Product `json:"products"`
	Meta     Meta             `json:"meta"`
	Message  string           `json:"message"`
}

type OneProductResponse struct {
	Product models.Product `json:"product"`
	Message string         `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ProductCreatePayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
}

type ProductUpdatePayload struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}
