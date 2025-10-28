package application

import (
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"
)

type GetProductService interface {
	GetProducts(page, size int) (dto.ProductDto, error)
}

type GetProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewGetProductService(productRepository repository.ProductRepository) GetProductService {

	return &GetProductServiceImpl{
		productRepository: productRepository,
	}
}

func (s *GetProductServiceImpl) GetProducts(page, size int) (dto.ProductDto, error) {

	return dto.ProductDto{}, nil
}
