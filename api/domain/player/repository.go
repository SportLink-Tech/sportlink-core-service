package player

type Repository interface {
	Insert(entity Entity) error
	Update(entity Entity) error
	Find(query DomainQuery) ([]Entity, error)
}

type DomainQuery interface {
}
