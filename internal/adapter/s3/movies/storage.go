package movies

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	errs "github.com/drobyshevv/movie-reservation/internal/errors"
	movieuc "github.com/drobyshevv/movie-reservation/internal/usecase/movies"
	"github.com/google/uuid"
)

const posterPrefix = "movies/posters/"

type S3MovieStorage struct {
	client *s3.Client
	bucket string
}

func NewS3MovieStorage(client *s3.Client, bucket string) *S3MovieStorage {
	return &S3MovieStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *S3MovieStorage) PutPoster(ctx context.Context, file io.ReadSeeker, size int64, contentType string) (uuid.UUID, error) {
	const op = "s3.movies.PutPoster"

	key := uuid.New()
	keyStr := fmt.Sprintf("%s%s", posterPrefix, key.String())

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &s.bucket,
		Key:           &keyStr,
		Body:          file,
		ContentLength: &size,
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return key, nil
}

func (s *S3MovieStorage) GetPoster(ctx context.Context, key uuid.UUID) (*movieuc.File, error) {
	const op = "s3.movies.GetPoster"

	keyStr := fmt.Sprintf("%s%s", posterPrefix, key.String())

	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &keyStr,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchKey" {
				return nil, fmt.Errorf("%s: %w", op, errs.ErrPosterNotFound)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movieuc.File{
		Body:        output.Body,
		ContentType: aws.ToString(output.ContentType),
		Size:        output.ContentLength,
	}, nil
}

func (s *S3MovieStorage) DeletePoster(ctx context.Context, key uuid.UUID) error {
	const op = "s3.movies.DeletePoster"

	keyStr := fmt.Sprintf("%s%s", posterPrefix, key.String())

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &keyStr,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
