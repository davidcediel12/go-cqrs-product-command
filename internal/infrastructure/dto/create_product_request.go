package dto

type CreateProductRequest struct {
	Name   string               `json:"name"`
	Price  float64              `json:"price"`
	Stock  float64              `json:"stock"`
	Images []CreateProductImage `json:"images"`
}

type CreateProductImage struct {
	Url       string `json:"url"`
	IsPrimary bool   `json:"isPrimary"`
}
