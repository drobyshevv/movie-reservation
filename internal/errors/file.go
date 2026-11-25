package errs

import "errors"

var (
	ErrEmptyFile         = errors.New("file is empty")
	ErrSmallFileSize     = errors.New("small file size")
	ErrInvalidPosterType = errors.New("invalid poster type")
	ErrOrphanedPoster    = errors.New("orphaned poster")
	ErrPosterNotFound    = errors.New("poster not found")
)
