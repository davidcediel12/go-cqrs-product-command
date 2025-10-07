package controller

import (
	"cqrs/command/internal/application"
	"cqrs/command/internal/infrastructure/dto"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	createProductService    application.CreateProductService
	generateImageUrlService application.GenerateImageUrlService
}

func NewProductController(
	createProductService application.CreateProductService,
	generateImageUrlService application.GenerateImageUrlService) *ProductController {

	return &ProductController{
		createProductService:    createProductService,
		generateImageUrlService: generateImageUrlService,
	}
}

func (c *ProductController) CreateProduct(ctx *fiber.Ctx) error {

	var productRequest dto.CreateProductRequest

	if err := ctx.BodyParser(&productRequest); err != nil {
		fmt.Println(err)
		return getInvalidRequest(ctx)
	}

	fmt.Println(productRequest)

	product, _ := c.createProductService.CreateProduct(ctx.UserContext(), &productRequest)

	return ctx.Status(fiber.StatusCreated).JSON(product)
}

func (c *ProductController) CreateImageUrls(ctx *fiber.Ctx) error {

	imageNames, err := getImageNames(ctx)

	if err != nil {
		return getInvalidRequest(ctx)
	}

	urls, err := c.generateImageUrlService.GenerateUrls(ctx.UserContext(), imageNames)

	if err != nil {

		fmt.Errorf("Error creating the urls: %w", err)

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	imageUrlResponse := dto.ImageUrlsResponse{
		Urls: urls,
	}

	return ctx.Status(fiber.StatusOK).JSON(imageUrlResponse)

}

func getInvalidRequest(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(
		fiber.Map{
			"error": "Invalid request body",
		},
	)
}

func getImageNames(ctx *fiber.Ctx) ([]string, error) {

	var imageUrlsRequest dto.ImageUrlsRequest

	if err := ctx.BodyParser(&imageUrlsRequest); err != nil {
		fmt.Println(err)

		return nil, getInvalidRequest(ctx)
	}

	imageNames := make([]string, 0, len(imageUrlsRequest.Images))

	for _, imageName := range imageUrlsRequest.Images {

		imageNames = append(imageNames, imageName.Name)
	}

	return imageNames, nil
}
