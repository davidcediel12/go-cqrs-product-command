package controller

import (
	"cqrs/command/internal/application"
	"cqrs/command/internal/infrastructure/dto"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	createProductService application.CreateProductService
}

func NewProductController(createProductService application.CreateProductService) *ProductController {
	return &ProductController{
		createProductService: createProductService,
	}
}

func (c *ProductController) CreateProduct(ctx *fiber.Ctx) error {

	var productRequest dto.CreateProductRequest

	if err := ctx.BodyParser(&productRequest); err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "Invalid request body",
			},
		)
	}

	fmt.Println(productRequest)

	_, product := c.createProductService.CreateProduct(ctx.UserContext(), &productRequest)

	return ctx.Status(fiber.StatusCreated).JSON(product)
}
