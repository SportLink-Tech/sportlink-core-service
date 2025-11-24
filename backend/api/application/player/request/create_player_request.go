package request

// NewPlayerRequest defines the structure of the request body for the player creation endpoint.
type NewPlayerRequest struct {
	Sport    string `json:"sport" validate:"required,oneof=Football Paddle Tennis"`
	ID       string `json:"id" validate:"required"`
	Category int    `json:"category" validate:"omitempty"`
}


