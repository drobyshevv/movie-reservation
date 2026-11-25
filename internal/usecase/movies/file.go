package movies

import "io"

type File struct {
	Body        io.ReadCloser
	ContentType string
	Size        *int64
}
