package request

// UpdateTeamRequest defines the body for the team update (PATCH) endpoint.
// All fields are optional; only provided fields are applied.
type UpdateTeamRequest struct {
	Name string `json:"name" validate:"omitempty"`
}
