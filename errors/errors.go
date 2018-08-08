package errors

import "fmt"

//CustomError is the base for all other CustomError types
type CustomError struct {
	Message string
	Code    int
}

func (e CustomError) Error() string {
	return e.Message
}

// NewCRUDError creates a new missing field error
func NewCRUDError(msg string) CustomError {
	return CustomError{msg, 200}
}

// NewMissingFieldError creates a new missing field error
func NewMissingFieldError(field string) CustomError {
	msg := fmt.Sprintf("'%s' field is required.", field)
	return CustomError{msg, 001}
}

// NewInvalidFieldError creates a new invalid field error
func NewInvalidFieldError(field string) CustomError {
	msg := fmt.Sprintf("Invalid field '%s'.", field)
	return CustomError{msg, 002}
}

// NewInvalidTokenError returns a new invalid token error
func NewInvalidTokenError() CustomError {
	return CustomError{"Invalid token.", 003}
}

// NewInvalidCredentialsError returns a new invalid credentials error
func NewInvalidCredentialsError() CustomError {
	return CustomError{"Invalid credentials.", 004}
}

// NewInternalError creates a new internal error
func NewInternalError(msg string) CustomError {
	return CustomError{msg, 101}
}

// NewGenericError creates a new internal error
func NewGenericError() CustomError {
	return CustomError{"Oops, something went wrong, please try again.", 102}
}
