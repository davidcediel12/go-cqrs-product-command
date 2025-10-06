package application

import (
	"cqrs/command/internal/infrastructure/dto"

	"github.com/google/uuid"
)

type CreateProductServiceImpl struct{}

func NewProductService() CreateProductService {
	return &CreateProductServiceImpl{}
}

func (s *CreateProductServiceImpl) CreateProduct(createProductRequest *dto.CreateProductRequest) dto.ProductDto {

	return dto.ProductDto{
		Id:    uuid.New(),
		Name:  createProductRequest.Name,
		Price: createProductRequest.Price,
	}
}
