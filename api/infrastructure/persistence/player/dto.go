package player

import (
	"fmt"
	"sportlink/api/domain/player"
)

type Dto struct {
	EntityId string `dynamodbav:"EntityId"`
	Id       string `dynamodbav:"Id"`
	Category int    `dynamodbav:"Category"`
	Sport    string `dynamodbav:"Sport"`
}

func From(entity player.Entity) (Dto, error) {
	if entity.ID == "" {
		return Dto{}, fmt.Errorf("ID could not be empty")
	}

	return Dto{
		EntityId: "Entity#Player",
		Id:       entity.ID,
		Category: int(entity.Category),
		Sport:    string(entity.Sport),
	}, nil
}
