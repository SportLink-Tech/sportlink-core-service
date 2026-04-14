package response

import "time"

type TimeSlotResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type MatchResponse struct {
	ID           string            `json:"id"`
	Participants []string          `json:"participants"`
	Sport        string            `json:"sport"`
	Day          time.Time         `json:"day"`
	TimeSlot     *TimeSlotResponse `json:"time_slot,omitempty"`
	Title        string            `json:"title,omitempty"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
}
