package player

import "sportlink/api/domain/common"

type Repository interface {
	Save(entity Entity) error
	Find(query DomainQuery) ([]Entity, error)
}

type DomainQuery struct {
	ID       string
	Category common.Category
	Sport    common.Sport
}

func NewDomainQuery(ID string, category common.Category, sport common.Sport) *DomainQuery {
	return &DomainQuery{ID: ID, Category: category, Sport: sport}
}
