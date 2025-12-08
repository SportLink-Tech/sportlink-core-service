package team

import (
	"context"
	"sportlink/api/domain/common"
)

type Repository interface {
	Save(ctx context.Context, entity Entity) error
	Find(ctx context.Context, query DomainQuery) ([]Entity, error)
}

type ID struct {
	Name  string
	Sport common.Sport
}

type DomainQuery struct {
	Name       string
	Ids        []string
	Categories []common.Category
	Sports     []common.Sport
}
