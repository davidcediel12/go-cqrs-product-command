package routes

import (
	"cqrs/command/internal/infrastructure/controller"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, productController controller.ProductController) {

	app.Post("/products", productController.CreateProduct)
	app.Post("/products/images/generate", productController.CreateImageUrls)
}
