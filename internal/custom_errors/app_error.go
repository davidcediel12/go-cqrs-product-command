package customerrors

import "fmt"

type ErrorType string

const (
	ValidationError ErrorType = "VALIDATION_ERROR"
	InternalError   ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func NewAppError(errorType ErrorType, message string, err error) error {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

func (a *AppError) Error() string {

	if a.Err != nil {
		return fmt.Sprintf("%s: %s - %v", a.Type, a.Message, a.Err.Error())
	}

	return fmt.Sprintf("%s: %s", a.Type, a.Message)
}

func (a *AppError) Unwrap() error {
	return a.Err
}
