package errors

import "fmt"

func InvalidRequestFormat() AppError {
	return AppError{
		Code:    InvalidRequestFormatErrorCode,
		Message: "Invalid request format provided.",
	}
}

func RequestValidationFailed(message string) AppError {
	return AppError{
		Code:    RequestValidationFailedErrorCode,
		Message: fmt.Sprintf("request validation failed. Err: %s", message),
	}
}
