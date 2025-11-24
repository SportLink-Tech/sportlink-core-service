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

func UseCaseExecutionFailed(message string) AppError {
	return AppError{
		Code:    UseCaseExecutionErrorCode,
		Message: fmt.Sprintf("use case execution failed. Err: %s", message),
	}
}

func NotFound(message string) AppError {
	return AppError{
		Code:    NotFoundErrorCode,
		Message: message,
	}
}
