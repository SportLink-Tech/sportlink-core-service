package player

import "sportlink/api/domain/common"

type Repository interface {
	Save(entity Entity) error
	Find(query DomainQuery) ([]Entity, error)
}

type DomainQuery struct {
	Id       string
	Ids      []string
	Category common.Category
	Sport    common.Sport
}

func NewDomainQuery(
	id string,
	ids []string,
	category common.Category,
	sport common.Sport,
) *DomainQuery {
	return &DomainQuery{
		Id:       id,
		Ids:      ids,
		Category: category,
		Sport:    sport,
	}
}
