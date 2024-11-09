package request

// NewTeamRequest defines the structure of the request body for the team creation endpoint.
type NewTeamRequest struct {
	Sport     string   `json:"sport" validate:"required,oneof=Football Paddle"`
	Name      string   `json:"name" validate:"required"`
	Category  int      `json:"category" validate:"omitempty"`
	PlayerIds []string `json:"players" validate:"omitempty"`
}
