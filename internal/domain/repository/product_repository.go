package repository

import "cqrs/command/internal/infrastructure/dto"

type ProductRepository interface {
	CreateProduct(createProductRequest *dto.CreateProductRequest) dto.ProductDto
}
