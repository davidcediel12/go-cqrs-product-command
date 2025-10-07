package application

import (
	"context"
	"cqrs/command/internal/infrastructure/dto"
)

type CreateProductService interface {
	CreateProduct(
		ctx context.Context,
		createProductRequest *dto.CreateProductRequest,
	) (dto.ProductDto, error)
}
