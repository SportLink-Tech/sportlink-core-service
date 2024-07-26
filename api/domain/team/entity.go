package team

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
)

type Entity struct {
	name     string
	category common.Category
	stats    common.Stats
	sport    common.Sport
	members  []player.Entity
}

func NewTeam(
	name string,
	category common.Category,
	stats common.Stats,
	sport common.Sport,
	members []player.Entity,
) *Entity {
	return &Entity{
		name:     name,
		category: category,
		stats:    stats,
		sport:    sport,
		members:  members,
	}
}
