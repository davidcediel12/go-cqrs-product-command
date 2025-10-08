package application

import (
	"context"
	"cqrs/command/internal/application/ports"
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"
	"cqrs/command/internal/logger"
)

type CreateProductServiceImpl struct {
	productRepository repository.ProductRepository
	messagePublisher  ports.MessagePublisher
}

func NewProductService(productRepository repository.ProductRepository,
	messagePublisher ports.MessagePublisher) CreateProductService {

	return &CreateProductServiceImpl{
		productRepository: productRepository,
		messagePublisher:  messagePublisher,
	}
}

func (s *CreateProductServiceImpl) CreateProduct(ctx context.Context,
	createProductRequest *dto.CreateProductRequest) (dto.ProductDto, error) {

	productDto, err := s.productRepository.CreateProduct(ctx, createProductRequest)

	if err != nil {
		return dto.ProductDto{}, err
	}

	logger.Log.Infof("Product %v created", productDto.Id)

	err = s.messagePublisher.PublishNewProduct(ctx, "product", &productDto) // TODO change to env var

	if err != nil {
		return dto.ProductDto{}, err
	}

	logger.Log.Infof("Message to indicate that product %v was created published", productDto.Id)

	return productDto, nil
}
