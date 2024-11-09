package errors

type ErrorCode string

const (
	InvalidRequestFormatErrorCode    ErrorCode = "invalid_request_format"
	RequestValidationFailedErrorCode ErrorCode = "request_validation_failed"
	NotFoundErrorCode                ErrorCode = "not_found"
	UnauthorizedErrorCode            ErrorCode = "unauthorized"
	UnexpectedErrorCode              ErrorCode = "unexpected_error"
	UseCaseExecutionErrorCode        ErrorCode = "use_case_execution_error"
)

type AppError struct {
	Code    ErrorCode
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
