package main

import (
	"context"
	"cqrs/command/internal/application"
	"cqrs/command/internal/infrastructure/controller"
	"cqrs/command/internal/infrastructure/messaging"
	"cqrs/command/internal/infrastructure/persistence"
	"cqrs/command/internal/infrastructure/routes"
	"cqrs/command/internal/infrastructure/storage"
	"cqrs/command/internal/logger"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	productController *controller.ProductController
)

func main() {

	pool := connectToPostgres()
	defer pool.Close()

	sdkConfig, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatal("Unable to connect to AWS services")
	}

	logger.Init()

	injectDependencies(sdkConfig, pool)

	app := fiber.New()

	routes.Routes(app, *productController)

	app.Listen((":8080"))
}

func connectToPostgres() *pgxpool.Pool {
	connStr := "postgres://postgres:admin@localhost:5432/products?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatal("Unable to connect to db: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")

	return pool

}

func injectDependencies(sdkConfig aws.Config, pool *pgxpool.Pool) {

	snsClient := sns.NewFromConfig(sdkConfig)

	s3Client := s3.NewFromConfig(sdkConfig)

	s3PresignClient := s3.NewPresignClient(s3Client)

	messagePublisher := messaging.NewSnsPublisher(snsClient)

	storageService := storage.NewStorageService(s3PresignClient)
	generateImageService := application.NewGenerateImageService(storageService)

	productRepository := persistence.NewProductRepository(pool)

	createProductService := application.NewProductService(productRepository, messagePublisher)

	productController = controller.NewProductController(createProductService, generateImageService)
}
