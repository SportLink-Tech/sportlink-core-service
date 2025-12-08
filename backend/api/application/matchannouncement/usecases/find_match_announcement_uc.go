package usecases

import (
	"context"
	"sportlink/api/domain/matchannouncement"
)

// FindMatchAnnouncementResult contains the paginated results and metadata
type FindMatchAnnouncementResult struct {
	Entities []matchannouncement.Entity
	Page     PageInfo
}

// PageInfo contains pagination metadata
type PageInfo struct {
	Number int // Current page number (1-based)
	OutOf  int // Total number of pages
	Total  int // Total number of items matching the query
}

type FindMatchAnnouncementUC struct {
	matchAnnouncementRepository matchannouncement.Repository
}

func NewFindMatchAnnouncementUC(matchAnnouncementRepository matchannouncement.Repository) *FindMatchAnnouncementUC {
	return &FindMatchAnnouncementUC{
		matchAnnouncementRepository: matchAnnouncementRepository,
	}
}

func (uc *FindMatchAnnouncementUC) Invoke(ctx context.Context, query matchannouncement.DomainQuery) (*FindMatchAnnouncementResult, error) {
	page, err := uc.matchAnnouncementRepository.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	pageInfo := CalculatePageInfo(query.Limit, query.Offset, page.Total)

	return &FindMatchAnnouncementResult{
		Entities: page.Entities,
		Page:     pageInfo,
	}, nil
}

// CalculatePageInfo calculates the current page number and total pages based on limit, offset, and total
func CalculatePageInfo(limit, offset, total int) PageInfo {
	var currentPage int
	var totalPages int

	if limit > 0 {
		// Page number is 1-based: offset=0 -> page 1, offset=limit -> page 2, etc.
		currentPage = (offset / limit) + 1
		// Total pages is ceiling of total / limit
		totalPages = (total + limit - 1) / limit
		if totalPages == 0 && total > 0 {
			totalPages = 1
		}
	} else {
		// No limit means all results in one page
		currentPage = 1
		if total > 0 {
			totalPages = 1
		} else {
			totalPages = 0
		}
	}

	return PageInfo{
		Number: currentPage,
		OutOf:  totalPages,
		Total:  total,
	}
}
