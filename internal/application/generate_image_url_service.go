package application

import (
	"context"
	"cqrs/command/internal/application/ports"
)

type GenerateImageUrlService interface {
	GenerateUrls(ctx context.Context, imageNames []string) ([]string, error)
}

type GenerateImageUrlServiceImpl struct {
	storageService ports.StorageService
}

func NewGenerateImageService(storageService ports.StorageService) GenerateImageUrlService {
	return &GenerateImageUrlServiceImpl{
		storageService: storageService,
	}
}

func (s *GenerateImageUrlServiceImpl) GenerateUrls(ctx context.Context, imageNames []string) ([]string, error) {

	var urls = make([]string, 0, len(imageNames))

	for _, imageName := range imageNames {

		url, err := s.storageService.GenerateUrl(ctx, imageName)

		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}
