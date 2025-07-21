package team

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/team"
)

// TODO faltan las stats
type Dto struct {
	EntityId string `dynamodbav:"EntityId"`
	Id       string `dynamodbav:"Id"`
	Category int    `dynamodbav:"Category"`
	Sport    string `dynamodbav:"Sport"`
}

func (d *Dto) ToDomain() team.Entity {
	return team.Entity{
		Name:     d.Id,
		Category: common.Category(d.Category),
		Sport:    common.Sport(d.Sport),
	}
}
