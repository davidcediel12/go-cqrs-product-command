package dto

import "github.com/google/uuid"

type ProductDto struct {
	Id     uuid.UUID            `json:"id"`
	Name   string               `json:"name"`
	Price  float64              `json:"price"`
	Stock  float64              `json:"stock"`
	Images []CreateProductImage `json:"images"`
}

type ProductImageDto struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	IsPrimary bool   `json:"isPrimary"`
}
