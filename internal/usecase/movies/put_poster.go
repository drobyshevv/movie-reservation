package movies

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	errs "github.com/drobyshevv/movie-reservation/internal/errors"
)

var allowedPosterTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// TODO реализовать удаление старого постера из s3
func (u *MovieUseCase) PutPoster(ctx context.Context, id int64, file io.ReadSeeker, size int64) error {
	const op = "usecase.movies.PutPoster"

	if size <= 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrEmptyFile)
	}

	var buf [512]byte

	n, err := io.ReadFull(file, buf[:])
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("%s: read header: %w", op, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrEmptyFile)
	}

	contentType := http.DetectContentType(buf[:n])
	if _, ok := allowedPosterTypes[contentType]; !ok {
		return fmt.Errorf("%s: %w: %s", op, errs.ErrInvalidPosterType, contentType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("%s: seek: %w", op, err)
	}

	posterID, err := u.storage.PutPoster(ctx, file, size, contentType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = u.repository.PutPosterID(ctx, id, posterID)
	if err != nil {
		cleanupCtx := context.WithoutCancel(ctx)
		if delErr := u.storage.DeletePoster(cleanupCtx, posterID); delErr != nil {
			return fmt.Errorf("%s: orphaned poster %s: %w (cleanup: %w)", op, posterID, err, delErr)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
