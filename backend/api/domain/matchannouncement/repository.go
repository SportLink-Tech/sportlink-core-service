package matchannouncement

import (
	"context"
	"sportlink/api/domain/common"
	"time"
)

// Page contains the results of a Find operation with pagination information
type Page struct {
	Entities []Entity // The matching entities
	Total    int      // Total number of entities matching the query (ignoring limit/offset)
}

// Repository defines the persistence operations for match announcements
type Repository interface {
	Save(ctx context.Context, entity Entity) error
	Find(ctx context.Context, query DomainQuery) (Page, error)
}

// DomainQuery represents the search criteria for match announcements
type DomainQuery struct {
	Sports     []common.Sport    // Search by multiple sports
	Categories []common.Category // Search by multiple admitted categories
	Statuses   []Status          // Search by multiple statuses
	FromDate   time.Time         // Announcements from this date
	ToDate     time.Time         // Announcements until this date
	Location   *Location         // Search by location (optional)
	Limit      int               // Maximum number of results to return (0 = no limit)
	Offset     int               // Number of results to skip (0 = no offset)
}

func NewDomainQuery(
	sports []common.Sport,
	categories []common.Category,
	statuses []Status,
	fromDate time.Time,
	toDate time.Time,
	location *Location,
	limit int,
	offset int,
) *DomainQuery {
	return &DomainQuery{
		Sports:     sports,
		Categories: categories,
		Statuses:   statuses,
		FromDate:   fromDate,
		ToDate:     toDate,
		Location:   location,
		Limit:      limit,
		Offset:     offset,
	}
}
