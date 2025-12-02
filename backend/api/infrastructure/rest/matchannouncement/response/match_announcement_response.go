package response

import "time"

// MatchAnnouncementResponse represents the API response for a match announcement
type MatchAnnouncementResponse struct {
	ID                 string                `json:"id,omitempty"`
	TeamName           string                `json:"team_name"`
	Sport              string                `json:"sport"`
	Day                time.Time             `json:"day"`
	TimeSlot           TimeSlotResponse      `json:"time_slot"`
	Location           LocationResponse      `json:"location"`
	AdmittedCategories CategoryRangeResponse `json:"admitted_categories"`
	Status             string                `json:"status"`
	CreatedAt          time.Time             `json:"created_at"`
}

type TimeSlotResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type LocationResponse struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	Locality string `json:"locality"`
}

type CategoryRangeResponse struct {
	Type       string `json:"type"`
	Categories []int  `json:"categories,omitempty"`
	MinLevel   int    `json:"min_level,omitempty"`
	MaxLevel   int    `json:"max_level,omitempty"`
}
