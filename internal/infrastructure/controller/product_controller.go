package controller

import (
	"cqrs/command/internal/application"
	customerrors "cqrs/command/internal/custom_errors"
	"cqrs/command/internal/infrastructure/dto"
	"cqrs/command/internal/logger"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	CREATE_PRODUCT      string = "create_product"
	GENERATE_IMAGE_URLS string = "generate_image_urls"
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
		logger.Log.WithError(err).Error("Error while parsing the body to create a product")
		return getInvalidRequest(ctx)
	}

	logger.Log.Infof("Received request to create product %v", productRequest.Name)

	product, err := c.createProductService.CreateProduct(ctx.UserContext(), &productRequest)

	if err != nil {
		return manageError(ctx, err, CREATE_PRODUCT)
	}

	return ctx.Status(fiber.StatusCreated).JSON(product)
}

func (c *ProductController) CreateImageUrls(ctx *fiber.Ctx) error {

	imageNames, err := getImageNames(ctx)

	if err != nil {
		return getInvalidRequest(ctx)
	}

	logger.Log.Infof("Received request to generate urls for %v", imageNames)

	urls, err := c.generateImageUrlService.GenerateUrls(ctx.UserContext(), imageNames)

	if err != nil {
		return manageError(ctx, err, GENERATE_IMAGE_URLS)
	}

	imageUrlResponse := dto.ImageUrlsResponse{
		Urls: urls,
	}

	return ctx.Status(fiber.StatusOK).JSON(imageUrlResponse)

}

func getInternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(
		fiber.Map{
			"error": "Internal server error",
		},
	)
}

func getInvalidRequest(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(
		fiber.Map{
			"error": "Invalid request",
		},
	)
}

func manageError(ctx *fiber.Ctx, err error, event string) error {

	var appErr *customerrors.AppError

	logger.Log.WithError(err).Errorf("Error during event %v ", event)

	if !errors.As(err, &appErr) {
		return getInternalServerError(ctx)
	}

	switch appErr.Type {
	case customerrors.ValidationError:

		return ctx.Status(fiber.ErrBadRequest.Code).JSON(
			fiber.Map{
				"error": appErr.Message,
			},
		)

	case customerrors.InternalError:

		return getInternalServerError(ctx)

	default:
		return getInternalServerError(ctx)
	}
}

func getImageNames(ctx *fiber.Ctx) ([]string, error) {

	var imageUrlsRequest dto.ImageUrlsRequest

	if err := ctx.BodyParser(&imageUrlsRequest); err != nil {
		logger.Log.WithError(err).Error("Error while parsing the body to generate image urls")

		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	imageNames := make([]string, 0, len(imageUrlsRequest.Images))

	for _, imageName := range imageUrlsRequest.Images {

		imageNames = append(imageNames, imageName.Name)
	}

	return imageNames, nil
}
