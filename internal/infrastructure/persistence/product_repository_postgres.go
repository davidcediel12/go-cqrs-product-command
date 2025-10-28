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

	logger.Log.Infof("Executing query to retrieve products with: %v", productsCommand)

	offset := productsCommand.Page * productsCommand.Size

	if productsCommand.Images {
		return r.getProductsWithImages(ctx, productsCommand, offset)
	}

	return r.getProducts(ctx, productsCommand, offset)

}

func (r *ProductRepositoryImpl) getProductsWithImages(ctx context.Context,
	productsCommand command.GetProductsCommand, offset int) ([]dto.ProductDto, error) {

	queryProductIds := `SELECT p.id from products p order by p.product_name offset $1 limit $2`

	rows, err := r.pool.Query(ctx, queryProductIds, offset, productsCommand.Size)

	if err != nil {
		return []dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"error while performing the query to obtain products", err)
	}

	var productIds []uuid.UUID

	for rows.Next() {

		var productId uuid.UUID

		rows.Scan(&productId)

		productIds = append(productIds, productId)
	}

	logger.Log.Infof("Obtained the products: %v to then obtain the images", productIds)

	queryGetProducts := `SELECT p.id, p.product_name, p.price, p.stock, 
	pi.id  as image_id, 
	COALESCE(pi.url, ''), 
	COALESCE(pi.is_primary, false)
	FROM products p left JOIN product_images pi ON p.id = pi.product_id
	WHERE p.id = ANY($1)
	ORDER BY p.product_name`

	productsMap := make(map[uuid.UUID]*dto.ProductDto)

	rows, err = r.pool.Query(ctx, queryGetProducts, productIds)

	if err != nil {
		return []dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"error while performing the query to obtain products", err)
	}

	defer rows.Close()

	for rows.Next() {
		var product productRow

		if err = rows.Scan(&product.ProductID, &product.ProductName, &product.Price, &product.Stock,
			&product.ImageID, &product.ImageURL, &product.IsPrimary); err != nil {

			return nil, customerrors.NewAppError(customerrors.InternalError,
				"error while mapping sql rows to product object", err)
		}

		productDto, exists := productsMap[product.ProductID]

		imageDto := dto.ProductImageDto{
			Id:        product.ImageID.String(),
			Url:       product.ImageURL,
			IsPrimary: product.IsPrimary,
		}

		if !exists {

			productsMap[product.ProductID] = &dto.ProductDto{
				Id:     product.ProductID,
				Name:   product.ProductName,
				Price:  product.Price,
				Stock:  product.Stock,
				Images: []dto.ProductImageDto{imageDto},
			}
		} else {
			productDto.Images = append(productDto.Images, imageDto)
		}
	}

	productsDto := make([]dto.ProductDto, 0, len(productsMap))

	for _, product := range productsMap {
		productsDto = append(productsDto, *product)
	}

	return productsDto, nil

}

func (r *ProductRepositoryImpl) getProducts(ctx context.Context, productsCommand command.GetProductsCommand,
	offset int) ([]dto.ProductDto, error) {

	queryGetProducts := `SELECT p.id, p.product_name, p.price, p.stock FROM products p 
	ORDER BY p.product_name offset $1 limit $2`

	rows, err := r.pool.Query(ctx, queryGetProducts, offset, productsCommand.Size)

	if err != nil {
		return []dto.ProductDto{}, customerrors.NewAppError(customerrors.InternalError,
			"error while performing the query to obtain products", err)
	}

	defer rows.Close()

	var products []dto.ProductDto

	for rows.Next() {

		var product dto.ProductDto

		if err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, customerrors.NewAppError(customerrors.InternalError,
				"error while mapping from sql row to object while obtaining products", err)
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

type productRow struct {
	ProductID   uuid.UUID
	ProductName string
	Price       float64
	Stock       float64
	ImageID     uuid.UUID
	ImageURL    string
	IsPrimary   bool
}
