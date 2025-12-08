package membership

import (
	"fmt"
	"sportlink/api/domain/common"
)

type Entity struct {
	ID       string
	TeamID   string
	PlayerID string
	Sport    common.Sport
}

func NewMembership(
	teamID string,
	playerID string,
	sport common.Sport) Entity {
	return Entity{
		ID:       generateTeamMembershipID(playerID, sport, teamID),
		TeamID:   teamID,
		PlayerID: playerID,
		Sport:    sport,
	}
}

func generateTeamMembershipID(playerId string, sport common.Sport, name string) string {
	return fmt.Sprintf("PLAYER#%s#SPORT#%s#TEAM#%s", playerId, sport, name)
}
