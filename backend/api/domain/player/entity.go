package player

import (
	"sportlink/api/domain/common"
)

type Entity struct {
	ID       string
	Category common.Category
	Sport    common.Sport
}

func NewPlayer(ID string, category common.Category) *Entity {
	return &Entity{ID: ID, Category: category}
}
