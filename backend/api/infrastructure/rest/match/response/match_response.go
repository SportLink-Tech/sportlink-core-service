package response

import "time"

type MatchResponse struct {
	ID           string    `json:"id"`
	Participants []string  `json:"participants"`
	Sport        string    `json:"sport"`
	Day          time.Time `json:"day"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
