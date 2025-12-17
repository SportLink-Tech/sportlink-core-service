package request

// NewPlayerRequest defines the structure of the request body for the player creation endpoint.
// ID is generated automatically using ULID, so it's not required in the request.
type NewPlayerRequest struct {
	Sport    string `json:"sport" validate:"required,oneof=Football Paddle Tennis"`
	Category int    `json:"category" validate:"omitempty"`
}
