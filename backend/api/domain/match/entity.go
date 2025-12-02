package match

import "time"

type Entity struct {
	ID                  string
	MatchAnnouncementID string
	HomeTeamID          string
	AwayTeamID          string
	MatchDate           time.Time
}
