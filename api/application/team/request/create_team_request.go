package request

// NewTeamRequest defines the structure of the request body for the team creation endpoint.
type NewTeamRequest struct {
	Sport    string   `json:"sport" validate:"required,oneof=football paddle"`
	Name     string   `json:"name" validate:"required"`
	Category int      `json:"category" validate:"omitempty"`
	Players  []string `json:"players" validate:"omitempty"`
}
