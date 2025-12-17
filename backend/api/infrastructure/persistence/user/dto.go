package user

import (
	"fmt"
	"sportlink/api/domain/user"
)

type Dto struct {
	EntityId  string   `dynamodbav:"EntityId"`
	Id        string   `dynamodbav:"Id"`
	FirstName string   `dynamodbav:"FirstName"`
	LastName  string   `dynamodbav:"LastName"`
	PlayerIDs []string `dynamodbav:"PlayerIDs"`
}

func From(entity user.Entity) (Dto, error) {
	if entity.ID == "" {
		return Dto{}, fmt.Errorf("id could not be empty")
	}

	return Dto{
		EntityId:  "Entity#User",
		Id:        fmt.Sprintf("ID#%s", entity.ID),
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		PlayerIDs: entity.PlayerIDs,
	}, nil
}

func (d *Dto) ToDomain() user.Entity {
	// Extract ID without the "ID#" prefix
	id := d.Id
	if len(id) > 3 && id[:3] == "ID#" {
		id = id[3:]
	}

	return user.Entity{
		ID:        id,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		PlayerIDs: d.PlayerIDs,
	}
}
