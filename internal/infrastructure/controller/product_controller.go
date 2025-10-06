package controller

import (
	"cqrs/command/internal/infrastructure/dto"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CreateProduct(ctx *fiber.Ctx) error {

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

	return ctx.SendString("Create Product")
}
