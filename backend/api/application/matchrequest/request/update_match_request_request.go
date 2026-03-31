package request

// UpdateMatchRequestRequest represents the HTTP request body for updating a match request status
type UpdateMatchRequestRequest struct {
	Status string `json:"status" validate:"required,oneof=ACCEPTED REJECTED"`
}
