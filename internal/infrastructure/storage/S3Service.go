package storage

import (
	"context"
	"cqrs/command/internal/application/ports"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	s3PresignClient *s3.PresignClient
}

func NewStorageService(s3PresignClient *s3.PresignClient) ports.StorageService {

	return &S3Service{
		s3PresignClient: s3PresignClient,
	}
}

func (s *S3Service) GenerateUrl(ctx context.Context, imageName string) (string, error) {

	request, err := s.s3PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("first-bucket-cqrs"), // TODO Change to ENV var
		Key:    aws.String("public/" + imageName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(5 * int64(time.Minute))
	})

	if err != nil {
		return "", fmt.Errorf("error while generating a signed URL for image %v, %w", imageName, err)
	}

	return request.URL, nil

}
