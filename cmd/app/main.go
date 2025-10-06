package main

import (
	"cqrs/command/internal/application"
	"cqrs/command/internal/infrastructure/controller"
	"cqrs/command/internal/infrastructure/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	createProductService := application.NewProductService()

	productController := controller.NewProductController(createProductService)

	app := fiber.New()

	routes.Routes(app, *productController)

	app.Listen((":8080"))
}
