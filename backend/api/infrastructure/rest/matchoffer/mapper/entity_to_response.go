package mapper

import (
	"sportlink/api/domain/matchoffer"
	"sportlink/api/infrastructure/rest/matchoffer/response"
)

func locationToResponse(loc matchoffer.Location) response.LocationResponse {
	r := response.LocationResponse{
		Country:  loc.Country,
		Province: loc.Province,
		Locality: loc.Locality,
	}
	if loc.HasCoords() {
		r.Latitude = &loc.Latitude
		r.Longitude = &loc.Longitude
	}
	return r
}

// EntityToResponse converts a domain entity to an API response DTO
func EntityToResponse(entity matchoffer.Entity) response.MatchOfferResponse {
	return response.MatchOfferResponse{
		ID:       entity.ID,
		TeamName: entity.TeamName,
		Sport:    string(entity.Sport),
		Day:      entity.Day,
		TimeSlot: response.TimeSlotResponse{
			StartTime: entity.TimeSlot.StartTime,
			EndTime:   entity.TimeSlot.EndTime,
		},
		Location:           locationToResponse(entity.Location),
		AdmittedCategories: categoryRangeToResponse(entity.AdmittedCategories),
		Status:             entity.Status.String(),
		CreatedAt:          entity.CreatedAt,
		OwnerAccountID:     entity.OwnerAccountID,
	}
}

// categoryRangeToResponse converts a domain CategoryRange to response DTO
func categoryRangeToResponse(cr matchoffer.CategoryRange) response.CategoryRangeResponse {
	resp := response.CategoryRangeResponse{
		Type: string(cr.Type),
	}

	switch cr.Type {
	case matchoffer.RangeTypeSpecific:
		categories := make([]int, len(cr.Categories))
		for i, cat := range cr.Categories {
			categories[i] = int(cat)
		}
		resp.Categories = categories
	case matchoffer.RangeTypeGreaterThan:
		resp.MinLevel = int(cr.MinLevel)
	case matchoffer.RangeTypeLessThan:
		resp.MaxLevel = int(cr.MaxLevel)
	case matchoffer.RangeTypeBetween:
		resp.MinLevel = int(cr.MinLevel)
		resp.MaxLevel = int(cr.MaxLevel)
	}

	return resp
}

// EntitiesToResponses converts a slice of domain entities to response DTOs
func EntitiesToResponses(entities []matchoffer.Entity) []response.MatchOfferResponse {
	responses := make([]response.MatchOfferResponse, len(entities))
	for i, entity := range entities {
		responses[i] = EntityToResponse(entity)
	}
	return responses
}

// NewPaginatedResponse creates a paginated response with the given data and pagination info
func NewPaginatedResponse(
	data []response.MatchOfferResponse,
	pageNumber int,
	totalPages int,
	total int,
) response.PaginatedMatchOffersResponse {
	return response.PaginatedMatchOffersResponse{
		Data: data,
		Pagination: response.PaginationInfo{
			Number: pageNumber,
			OutOf:  totalPages,
			Total:  total,
		},
	}
}
