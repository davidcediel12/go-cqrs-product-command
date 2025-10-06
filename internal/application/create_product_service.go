package application

import "cqrs/command/internal/infrastructure/dto"

type CreateProductService interface {
	CreateProduct(createProductRequest *dto.CreateProductRequest) dto.ProductDto
}
