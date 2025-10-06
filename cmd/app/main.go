package main

import (
	"cqrs/command/internal/application"
	"cqrs/command/internal/infrastructure/controller"
	"cqrs/command/internal/infrastructure/persistence"
	"cqrs/command/internal/infrastructure/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	productRepository := persistence.NewProductRepository()

	createProductService := application.NewProductService(productRepository)

	productController := controller.NewProductController(createProductService)

	app := fiber.New()

	routes.Routes(app, *productController)

	app.Listen((":8080"))
}
