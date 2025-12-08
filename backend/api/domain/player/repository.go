package player

import (
	"context"
	"sportlink/api/domain/common"
)

type Repository interface {
	Save(ctx context.Context, entity Entity) error
	Find(ctx context.Context, query DomainQuery) ([]Entity, error)
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
