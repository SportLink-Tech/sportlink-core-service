package account

import "context"

type Repository interface {
	Save(ctx context.Context, entity Entity) error
	Find(ctx context.Context, query DomainQuery) ([]Entity, error)
}

type DomainQuery struct {
	Id        string
	Ids       []string
	Email     string
	Emails    []string
	Nickname  string
	Nicknames []string
}
