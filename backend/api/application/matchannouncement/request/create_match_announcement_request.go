package request

// NewMatchAnnouncementRequest defines the structure of the request body for the match announcement creation endpoint.
type NewMatchAnnouncementRequest struct {
	TeamName           string             `json:"team_name" validate:"required"`
	Sport              string             `json:"sport" validate:"required"`
	Day                string             `json:"day" validate:"required"` // ISO date string
	TimeSlot           TimeSlot           `json:"time_slot" validate:"required"`
	Location           Location           `json:"location" validate:"required"`
	AdmittedCategories CategoryRangeInput `json:"admitted_categories" validate:"required"`
}

type TimeSlot struct {
	StartTime string `json:"start_time" validate:"required"` // ISO datetime string
	EndTime   string `json:"end_time" validate:"required"`   // ISO datetime string
}

type Location struct {
	Country  string `json:"country" validate:"required"`
	Province string `json:"province" validate:"required"`
	Locality string `json:"locality" validate:"required"`
}

type CategoryRangeInput struct {
	Type       string `json:"type" validate:"required,oneof=SPECIFIC GREATER_THAN LESS_THAN BETWEEN"`
	Categories []int  `json:"categories" validate:"omitempty"`
	MinLevel   int    `json:"min_level" validate:"omitempty"`
	MaxLevel   int    `json:"max_level" validate:"omitempty"`
}
