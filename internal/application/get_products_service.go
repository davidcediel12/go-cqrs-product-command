package application

import (
	"context"
	"cqrs/command/internal/command"
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"
)

type GetProductService interface {
	GetProducts(ctx context.Context, productsCommand command.GetProductsCommand) ([]dto.ProductDto, error)
}

type GetProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewGetProductService(productRepository repository.ProductRepository) GetProductService {

	return &GetProductServiceImpl{
		productRepository: productRepository,
	}
}

func (s *GetProductServiceImpl) GetProducts(ctx context.Context, productsCommand command.GetProductsCommand) ([]dto.ProductDto, error) {

	return s.productRepository.GetProducts(ctx, productsCommand)
}
