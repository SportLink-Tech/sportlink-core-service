package mapper

import (
	"sportlink/api/domain/matchannouncement"
	"sportlink/api/infrastructure/rest/matchannouncement/response"
)

// EntityToResponse converts a domain entity to an API response DTO
func EntityToResponse(entity matchannouncement.Entity) response.MatchAnnouncementResponse {
	return response.MatchAnnouncementResponse{
		ID:       entity.ID,
		TeamName: entity.TeamName,
		Sport:    string(entity.Sport),
		Day:      entity.Day,
		TimeSlot: response.TimeSlotResponse{
			StartTime: entity.TimeSlot.StartTime,
			EndTime:   entity.TimeSlot.EndTime,
		},
		Location: response.LocationResponse{
			Country:  entity.Location.Country,
			Province: entity.Location.Province,
			Locality: entity.Location.Locality,
		},
		AdmittedCategories: categoryRangeToResponse(entity.AdmittedCategories),
		Status:             entity.Status.String(),
		CreatedAt:          entity.CreatedAt,
	}
}

// categoryRangeToResponse converts a domain CategoryRange to response DTO
func categoryRangeToResponse(cr matchannouncement.CategoryRange) response.CategoryRangeResponse {
	resp := response.CategoryRangeResponse{
		Type: string(cr.Type),
	}

	switch cr.Type {
	case matchannouncement.RangeTypeSpecific:
		categories := make([]int, len(cr.Categories))
		for i, cat := range cr.Categories {
			categories[i] = int(cat)
		}
		resp.Categories = categories
	case matchannouncement.RangeTypeGreaterThan:
		resp.MinLevel = int(cr.MinLevel)
	case matchannouncement.RangeTypeLessThan:
		resp.MaxLevel = int(cr.MaxLevel)
	case matchannouncement.RangeTypeBetween:
		resp.MinLevel = int(cr.MinLevel)
		resp.MaxLevel = int(cr.MaxLevel)
	}

	return resp
}

// EntitiesToResponses converts a slice of domain entities to response DTOs
func EntitiesToResponses(entities []matchannouncement.Entity) []response.MatchAnnouncementResponse {
	responses := make([]response.MatchAnnouncementResponse, len(entities))
	for i, entity := range entities {
		responses[i] = EntityToResponse(entity)
	}
	return responses
}
