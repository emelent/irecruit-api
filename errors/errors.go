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

// CRUD creates a new CRUD error
func CRUD(msg string) CustomError {
	return CustomError{msg, 200}
}

// MissingField creates a new missing field error
func MissingField(field string) CustomError {
	msg := fmt.Sprintf("'%s' field is required.", field)
	return CustomError{msg, 001}
}

// InvalidField creates a new invalid field error
func InvalidField(field string) CustomError {
	msg := fmt.Sprintf("Invalid field '%s'.", field)
	return CustomError{msg, 002}
}

// InvalidToken returns a new invalid token error
func InvalidToken() CustomError {
	return CustomError{"Invalid token.", 003}
}

// InvalidCredentials returns a new invalid credentials error
func InvalidCredentials() CustomError {
	return CustomError{"Invalid credentials.", 004}
}

// Input returns a new invalid input error
func Input(msg string) CustomError {
	return CustomError{msg, 005}
}

// Internal creates a new internal error
func Internal(msg string) CustomError {
	return CustomError{msg, 101}
}

// Generic creates a new internal error
func Generic() CustomError {
	return CustomError{"Oops, something went wrong, please try again.", 102}
}
