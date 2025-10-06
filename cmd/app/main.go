package main

import (
	"context"
	"cqrs/command/internal/application"
	"cqrs/command/internal/infrastructure/controller"
	"cqrs/command/internal/infrastructure/persistence"
	"cqrs/command/internal/infrastructure/routes"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	pool := connectToPostgres()
	defer pool.Close()

	productRepository := persistence.NewProductRepository(pool)

	createProductService := application.NewProductService(productRepository)

	productController := controller.NewProductController(createProductService)

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
