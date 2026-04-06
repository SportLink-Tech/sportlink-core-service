package response

import "time"

type MatchRequestResponse struct {
	ID                  string    `json:"id"`
	MatchOfferID string    `json:"match_offer_id"`
	OwnerAccountID      string    `json:"owner_account_id"`
	RequesterAccountID  string    `json:"requester_account_id"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}
