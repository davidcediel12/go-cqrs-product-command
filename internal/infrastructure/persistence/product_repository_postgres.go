package persistence

import (
	"context"
	"cqrs/command/internal/command"
	customerrors "cqrs/command/internal/custom_errors"
	"cqrs/command/internal/domain/repository"
	"cqrs/command/internal/infrastructure/dto"
	"cqrs/command/internal/logger"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context,
	createProductRequest *dto.CreateProductRequest) (dto.ProductDto, error) {

	transaction, err := r.pool.Begin(ctx)

	if err != nil {
		return dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"error while starting the transaction", err)
	}

	defer func() {
		_ = transaction.Rollback(ctx)
	}()

	query := `
		INSERT INTO products(id, product_name, price, stock) values ($1, $2, $3, $4)
	`

	productId := uuid.New()

	_, err = transaction.Exec(ctx, query, productId, createProductRequest.Name,
		createProductRequest.Price, createProductRequest.Stock)

	if err != nil {
		return dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"saving product failed", err)
	}

	productImages, err := r.saveProductImages(ctx, transaction, createProductRequest.Images, productId)

	if err != nil {
		return dto.ProductDto{}, fmt.Errorf("saving product failed: %w", err)
	}

	if err := transaction.Commit(ctx); err != nil {
		return dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"error while commiting the transaction", err)
	}

	log.Infof("Product successfully saved with ID: %v", productId)

	return dto.ProductDto{
		Id:     productId,
		Name:   createProductRequest.Name,
		Price:  createProductRequest.Price,
		Stock:  createProductRequest.Stock,
		Images: productImages,
	}, nil
}

func (r *ProductRepositoryImpl) GetProducts(ctx context.Context, productsCommand command.GetProductsCommand) ([]dto.ProductDto, error) {

	logger.Log.Infof("Executing query to retrieve products with page %v and size %v", productsCommand.Page, productsCommand.Size)

	queryGetProducts := `SELECT p.id, p.product_name, p.price, p.stock FROM products p 
	ORDER BY p.product_name offset $1 limit $2`

	rows, err := r.pool.Query(ctx, queryGetProducts, productsCommand.Page, productsCommand.Size)

	if err != nil {
		return []dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"error while performing the query to obtain products", err)
	}

	defer rows.Close()

	var products []dto.ProductDto

	for rows.Next() {

		var product dto.ProductDto

		if err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock); err != nil {
			return []dto.ProductDto{}, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) saveProductImages(ctx context.Context, transaction pgx.Tx,
	productImages []dto.CreateProductImage, productId uuid.UUID) ([]dto.ProductImageDto, error) {

	createdImages := make([]dto.ProductImageDto, 0, len(productImages))
	queryProductImage := `
		INSERT INTO product_images(id, product_id, url, is_primary) values ($1, $2, $3, $4)
	`

	for _, image := range productImages {

		productImageId := uuid.New()

		_, err := transaction.Exec(ctx, queryProductImage, productImageId, productId,
			image.Url, image.IsPrimary)

		if err != nil {
			return []dto.ProductImageDto{}, fmt.Errorf("saving product image %s for product %s failed: %w", image.Url, productId, err)
		}

		createdImages = append(createdImages, dto.ProductImageDto{
			Id:        productImageId.String(),
			Url:       image.Url,
			IsPrimary: image.IsPrimary,
		})

	}

	return createdImages, nil

}
