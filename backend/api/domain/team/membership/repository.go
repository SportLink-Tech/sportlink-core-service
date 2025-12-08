package membership

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
	Name      string
	PlayerIDs []string
	Sports    []common.Sport
}
