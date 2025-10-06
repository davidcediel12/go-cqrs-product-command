package persistence

import (
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) repository.ProductRepository {
	return &ProductRepositoryImpl{
		pool: pool,
	}
}

func (r *ProductRepositoryImpl) CreateProduct(createProductRequest *dto.CreateProductRequest) dto.ProductDto {
	return dto.ProductDto{
		Id:    uuid.New(),
		Name:  "first product",
		Price: 182,
		Stock: 2092,
		Images: []dto.ProductImageDto{
			{
				Id:        uuid.New().String(),
				Url:       "my.url",
				IsPrimary: true,
			},
		},
	}
}
