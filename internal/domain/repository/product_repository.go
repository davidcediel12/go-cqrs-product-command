package repository

import (
	"context"
	"cqrs/command/internal/infrastructure/dto"
)

type ProductRepository interface {
	CreateProduct(
		context context.Context,
		createProductRequest *dto.CreateProductRequest,
	) (dto.ProductDto, error)

	GetProducts(page, size int) ([]dto.ProductDto, error)
}
