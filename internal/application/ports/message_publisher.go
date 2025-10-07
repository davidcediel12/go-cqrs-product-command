package ports

import (
	"context"
	"cqrs/command/internal/infrastructure/dto"
)

type MessagePublisher interface {
	PublishNewProduct(ctx context.Context, topic string, productDto *dto.ProductDto) error
}
