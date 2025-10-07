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

func (s *CreateProductServiceImpl) CreateProduct(createProductRequest *dto.CreateProductRequest, ctx context.Context) dto.ProductDto {

	return s.productRepository.CreateProduct(createProductRequest, ctx)
}
