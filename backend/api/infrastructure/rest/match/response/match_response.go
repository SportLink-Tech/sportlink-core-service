package response

import "time"

type MatchResponse struct {
	ID               string    `json:"id"`
	LocalAccountID   string    `json:"local_account_id"`
	VisitorAccountID string    `json:"visitor_account_id"`
	Sport            string    `json:"sport"`
	Day              time.Time `json:"day"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}
