package application

import (
	"context"
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"
)

type CreateProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) CreateProductService {

	return &CreateProductServiceImpl{
		productRepository: productRepository,
	}
}

func (s *CreateProductServiceImpl) CreateProduct(ctx context.Context,
	createProductRequest *dto.CreateProductRequest) (dto.ProductDto, error) {

	return s.productRepository.CreateProduct(ctx, createProductRequest)
}
