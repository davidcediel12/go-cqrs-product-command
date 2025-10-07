package ports

import "context"

type StorageService interface {
	GenerateUrl(ctx context.Context, imageName string) (string, error)
}
