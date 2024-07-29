package player

import (
	"sportlink/api/domain/player"
)

type Dto struct {
	EntityId string `dynamodbav:"EntityId"`
	Id       string `dynamodbav:"Id"`
	Category string `dynamodbav:"Category"`
	Sport    string `dynamodbav:"Sport"`
}

func From(entity player.Entity) Dto {
	return Dto{
		EntityId: "Entity#Player",
		Id:       entity.ID,
		Category: entity.Category.String(),
		Sport:    string(entity.Sport),
	}
}
