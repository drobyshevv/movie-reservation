package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	sysconfig "github.com/drobyshevv/movie-reservation/internal/config"
)

func New(cfg sysconfig.Config) (*s3.Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.S3.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.S3.AccessKey,
				cfg.S3.SecretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		o.UsePathStyle = true
	})

	return client, nil
}
