package errors

type ErrorCode string

const (
	InvalidRequestDataErrorCode ErrorCode = "invalid_request_data"
	NotFoundErrorCode           ErrorCode = "not_found"
	UnauthorizedErrorCode       ErrorCode = "unauthorized"
	UnexpectedErrorCode         ErrorCode = "unexpected_error"
)

type AppError struct {
	Code    ErrorCode
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
