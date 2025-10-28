package repository

import (
	"context"
	"cqrs/command/internal/command"
	"cqrs/command/internal/infrastructure/dto"
)

type ProductRepository interface {
	CreateProduct(
		context context.Context,
		createProductRequest *dto.CreateProductRequest,
	) (dto.ProductDto, error)

	GetProducts(ctx context.Context, productsCommand command.GetProductsCommand) ([]dto.ProductDto, error)
}
