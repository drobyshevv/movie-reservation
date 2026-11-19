package movies

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	movieuc "github.com/drobyshevv/movie-reservation/internal/usecase/movies"
	"github.com/google/uuid"
)

const posterPrefix = "movies/posters/"

type MediaStorage struct {
	client *s3.Client
	backet string
}

func NewPosterStoragey(client *s3.Client, backet string) *MediaStorage {
	return &MediaStorage{
		client: client,
		backet: backet,
	}
}

func (r *MediaStorage) PutPoster(ctx context.Context, file io.ReadSeeker, size int64) (uuid.UUID, error) {
	const op = "s3.movies.PutPoster"

	key := uuid.New()
	keyStr := fmt.Sprintf("%s%s.jpg", posterPrefix, key.String())

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &r.backet,
		Key:           &keyStr,
		Body:          file,
		ContentLength: &size,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return key, nil
}

func (r *MediaStorage) GetPoster(ctx context.Context, key uuid.UUID) (*movieuc.File, error) {
	const op = "s3.movies.GetPoster"

	keyStr := fmt.Sprintf("%s%s.jpg", posterPrefix, key.String())

	output, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &r.backet,
		Key:    &keyStr,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movieuc.File{
		Body:        output.Body,
		ContentType: aws.ToString(output.ContentType),
		Size:        output.ContentLength,
	}, nil
}

func (r *MediaStorage) DeletePoster(ctx context.Context, key uuid.UUID) error {
	const op = "s3.movies.DeletePoster"

	keyStr := fmt.Sprintf("%s%s.jpg", posterPrefix, key.String())

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &r.backet,
		Key:    &keyStr,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
