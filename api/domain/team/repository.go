package team

type Repository interface {
	Save(entity Entity) error
	Find(query DomainQuery) ([]Entity, error)
}

type DomainQuery struct {
	ID []string
}

func NewDomainQuery(ID []string) *DomainQuery {
	return &DomainQuery{ID: ID}

}
