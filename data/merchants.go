package data

import "time"

type Merchant struct {
	Id        string
	Name      string
	Email     string
	Password  string
	SkuId     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	SkuId    string `json:"skuId" validate:"required"`
}

type MerchantLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var MerchantData = []Merchant{}
