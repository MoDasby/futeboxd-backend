package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/modasby/futeboxd-backend/core/config"
)

func LoadConfig(cfg *config.AWS) (*aws.Config, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.Key, cfg.Secret, ""),
		),
		awsConfig.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, err
	}

	return &awsCfg, nil
}

func NewS3Client(cfg *config.AWS) (*s3.Client, error) {
	awsCfg, err := LoadConfig(cfg)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(*awsCfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		}
		o.UsePathStyle = true
	})

	return client, nil
}
