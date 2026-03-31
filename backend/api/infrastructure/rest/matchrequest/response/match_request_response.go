package response

import "time"

type MatchRequestResponse struct {
	ID                  string    `json:"id"`
	MatchAnnouncementID string    `json:"match_announcement_id"`
	OwnerAccountID      string    `json:"owner_account_id"`
	RequesterAccountID  string    `json:"requester_account_id"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}
