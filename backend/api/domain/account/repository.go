package account

import "context"

type Repository interface {
	Save(ctx context.Context, entity Entity) error
	Find(ctx context.Context, query DomainQuery) ([]Entity, error)
}

type DomainQuery struct {
	Ids       []string
	Emails    []string
	Nicknames []string
}
