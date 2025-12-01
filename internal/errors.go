package internal

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrDuplicateTitle     = errors.New("duplicate title")
	ErrTitleRequired      = errors.New("title is required")
	ErrTitleEmpty         = errors.New("title cannot be empty")
	ErrTitleMaxLength     = errors.New("title must be less than 255 characters")
	ErrInvalidID          = errors.New("id not valid")
	ErrEmptyList          = errors.New("list cannot be empty")
	ErrDuplicateInRequest = errors.New("duplicate entry in request")
	ErrLimitExceeded      = errors.New("limit exceeds maximum allowed")
)
