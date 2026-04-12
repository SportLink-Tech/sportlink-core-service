package events

// MatchOfferCapacityReachedEvent is published when all spots in a match offer
// have been accepted. A consumer can then call ConfirmMatchOfferUC to create
// the match automatically.
type MatchOfferCapacityReachedEvent struct {
	MatchOfferID   string
	OwnerAccountID string
}
