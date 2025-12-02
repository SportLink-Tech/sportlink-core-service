package matchannouncement

import (
	"sportlink/api/domain/common"
	"time"
)

// Repository defines the persistence operations for match announcements
type Repository interface {
	Save(entity Entity) error
	Find(query DomainQuery) ([]Entity, error)
}

// DomainQuery represents the search criteria for match announcements
type DomainQuery struct {
	Sports     []common.Sport    // Search by multiple sports
	Categories []common.Category // Search by multiple admitted categories
	Statuses   []Status          // Search by multiple statuses
	FromDate   time.Time         // Announcements from this date
	ToDate     time.Time         // Announcements until this date
	Location   *Location         // Search by location (optional)
}

func NewDomainQuery(
	sports []common.Sport,
	categories []common.Category,
	statuses []Status,
	fromDate time.Time,
	toDate time.Time,
	location *Location,
) *DomainQuery {
	return &DomainQuery{
		Sports:     sports,
		Categories: categories,
		Statuses:   statuses,
		FromDate:   fromDate,
		ToDate:     toDate,
		Location:   location,
	}
}
