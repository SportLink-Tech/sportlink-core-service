package team

import "sportlink/api/domain/common"

type Repository interface {
	Save(entity Entity) error
	Find(query DomainQuery) ([]Entity, error)
}

type DomainQuery struct {
	Name       string
	Ids        []string
	Categories []common.Category
}
