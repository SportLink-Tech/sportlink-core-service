package team

import (
	"fmt"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
)

type Entity struct {
	ID       string
	Name     string
	Category common.Category
	Stats    common.Stats
	Sport    common.Sport
	Members  []player.Entity
}

func NewTeam(
	name string,
	category common.Category,
	stats common.Stats,
	sport common.Sport,
	members []player.Entity,
) Entity {
	return Entity{
		ID:       generateTeamID(sport, name),
		Name:     name,
		Category: category,
		Stats:    stats,
		Sport:    sport,
		Members:  members,
	}
}

// generateTeamID creates a team ID in the format: SPORT#<sport>#NAME#<name>
func generateTeamID(sport common.Sport, name string) string {
	return fmt.Sprintf("SPORT#%s#NAME#%s", sport, name)
}
