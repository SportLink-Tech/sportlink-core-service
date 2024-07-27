package errors

func InvalidRequestData() AppError {
	return AppError{
		Code:    InvalidRequestDataErrorCode,
		Message: "Invalid request data provided",
	}
}
