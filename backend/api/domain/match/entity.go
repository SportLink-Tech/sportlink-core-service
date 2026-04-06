package match

import "time"

type Entity struct {
	ID                  string
	MatchOfferID string
	HomeTeamID          string
	AwayTeamID          string
	MatchDate           time.Time
}
