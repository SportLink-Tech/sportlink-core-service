package team

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
)

type Entity struct {
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
) *Entity {
	return &Entity{
		Name:     name,
		Category: category,
		Stats:    stats,
		Sport:    sport,
		Members:  members,
	}
}
