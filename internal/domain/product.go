package domain

import "github.com/google/uuid"

type Product struct {
	Id     uuid.UUID
	Name   string
	Price  float64
	Stock  int
	Images []ProductImage
}

type ProductImage struct {
	Id        uuid.UUID
	Url       string
	IsPrimary bool
}
