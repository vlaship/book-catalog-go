package apperr

// AppError is a struct for problem detail
type AppError struct {
	Err    error
	Title  string
	Code   string
	Detail string
}

// WithFunc applies a list of functions to an AppError and returns the modified AppError
func (p AppError) WithFunc(fn ...func(p AppError) AppError) AppError {
	p.Err = p
	for _, f := range fn {
		p = f(p)
	}
	return p
}

// WithTitle func
func WithTitle(title string) func(p AppError) AppError {
	return func(p AppError) AppError {
		p.Title = title
		return p
	}
}

// WithDetail func
func WithDetail(detail string) func(p AppError) AppError {
	return func(p AppError) AppError {
		p.Detail = detail
		return p
	}
}

// Error implement error interface
func (p AppError) Error() string {
	return p.Detail
}

// Unwrap implement error interface
func (p AppError) Unwrap() error {
	return p.Err
}
