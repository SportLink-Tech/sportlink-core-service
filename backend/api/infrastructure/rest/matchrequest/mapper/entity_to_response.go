package mapper

import (
	"sportlink/api/domain/matchrequest"
	"sportlink/api/infrastructure/rest/matchrequest/response"
)

func EntityToResponse(entity matchrequest.Entity) response.MatchRequestResponse {
	return response.MatchRequestResponse{
		ID:                  entity.ID,
		MatchAnnouncementID: entity.MatchAnnouncementID,
		OwnerAccountID:      entity.OwnerAccountID,
		RequesterAccountID:  entity.RequesterAccountID,
		Status:              entity.Status.String(),
		CreatedAt:           entity.CreatedAt,
	}
}

func EntitiesToResponses(entities []matchrequest.Entity) []response.MatchRequestResponse {
	responses := make([]response.MatchRequestResponse, len(entities))
	for i, e := range entities {
		responses[i] = EntityToResponse(e)
	}
	return responses
}
