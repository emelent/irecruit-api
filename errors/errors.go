package errors

import "fmt"

//CustomError is the base for all other CustomError types
type CustomError struct {
	Message string
}

//MissingFieldError is an error for missing fields
type MissingFieldError CustomError

func (e MissingFieldError) Error() string {
	return fmt.Sprintf("'%s' field is required.", e.Message)
}

//InvalidFieldError is an error for invalid fields
type InvalidFieldError CustomError

func (e InvalidFieldError) Error() string {
	return fmt.Sprintf("Invalid field '%s'.", e.Message)
}

//InvalidTokenError  error for invalid token
type InvalidTokenError CustomError

func (e InvalidTokenError) Error() string {
	return "Invalid token."
}

// NewMissingFieldError creates a new missing field error
func NewMissingFieldError(field string) MissingFieldError {
	return MissingFieldError{field}
}

// NewInvalidFieldError creates a new invalid field error
func NewInvalidFieldError(field string) InvalidFieldError {
	return InvalidFieldError{field}
}
