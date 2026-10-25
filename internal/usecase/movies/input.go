package movies

type CreateMovieInput struct {
	Title       string
	Description *string
	Duration    int32
}

type UpdateMovieInput struct {
	Title       *string
	Description *string
	Duration    *int32
}
