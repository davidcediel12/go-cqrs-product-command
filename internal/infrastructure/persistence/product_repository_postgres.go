package persistence

import (
	"context"
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"
	"log"

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

func (r *ProductRepositoryImpl) CreateProduct(createProductRequest *dto.CreateProductRequest, context context.Context) dto.ProductDto {

	query := `
		INSERT INTO products(id, product_name, price, stock) values ($1, $2, $3, $4)
	`

	productId := uuid.New()

	_, err := r.pool.Exec(context, query, productId, createProductRequest.Name,
		createProductRequest.Price, createProductRequest.Stock)

	if err != nil {
		log.Fatal("Saving product failed:", err)
	}

	productImages := r.saveProductImages(context, createProductRequest.Images, productId)

	return dto.ProductDto{
		Id:     productId,
		Name:   createProductRequest.Name,
		Price:  createProductRequest.Price,
		Stock:  createProductRequest.Stock,
		Images: productImages,
	}
}

func (r *ProductRepositoryImpl) saveProductImages(context context.Context,
	productImages []dto.CreateProductImage, productId uuid.UUID) []dto.ProductImageDto {

	var createdImages []dto.ProductImageDto

	queryProductImage := `
		INSERT INTO product_images(id, product_id, url, is_primary) values ($1, $2, $3, $4)
	`

	for _, image := range productImages {

		productImageId := uuid.New()

		_, err := r.pool.Exec(context, queryProductImage, productImageId, productId,
			image.Url, image.IsPrimary)

		if err != nil {
			log.Fatal("Saving product image failed:", err)
		}

		createdImages = append(createdImages, dto.ProductImageDto{
			Id:        productImageId.String(),
			Url:       image.Url,
			IsPrimary: image.IsPrimary,
		})

	}

	return createdImages

}
