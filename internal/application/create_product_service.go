package application

import (
	"context"
	"cqrs/command/internal/infrastructure/dto"
)

type CreateProductService interface {
	CreateProduct(createProductRequest *dto.CreateProductRequest, context context.Context) dto.ProductDto
}
