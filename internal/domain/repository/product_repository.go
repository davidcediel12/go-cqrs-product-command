package repository

import (
	"context"
	"cqrs/command/internal/infrastructure/dto"
)

type ProductRepository interface {
	CreateProduct(createProductRequest *dto.CreateProductRequest, context context.Context) (dto.ProductDto, error)
}
