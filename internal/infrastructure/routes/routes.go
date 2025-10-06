package routes

import (
	"cqrs/command/internal/infrastructure/controller"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {

	app.Post("/products", controller.CreateProduct)
}
