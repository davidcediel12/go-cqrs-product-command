package messaging

import (
	"context"
	"cqrs/command/internal/application/ports"
	"cqrs/command/internal/infrastructure/dto"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsPublisher struct {
	snsClient *sns.Client
}

func NewSnsPublisher(snsClient *sns.Client) ports.MessagePublisher {

	return &SnsPublisher{
		snsClient: snsClient,
	}
}

func (p *SnsPublisher) PublishNewProduct(ctx context.Context, topic string, productDto *dto.ProductDto) error {

	jsonData, err := json.Marshal(productDto)

	if err != nil {
		return fmt.Errorf("error while converting the product %v into json: %w", productDto, err)
	}

	publishInput := sns.PublishInput{
		TopicArn: aws.String("arn:aws:sns:us-east-2:641675857246:" + topic), // TODO Move to ENV variables
		Message:  aws.String(string(jsonData)),
	}

	_, err = p.snsClient.Publish(ctx, &publishInput)

	if err != nil {
		log.Printf("cannot publish message %v to topic %v: %v", string(jsonData), topic, err)

		return fmt.Errorf("cannot publish message %v to topic %v: %w", string(jsonData), topic, err)
	}

	log.Printf("product %v published correctly", productDto.Id)

	return err

}
