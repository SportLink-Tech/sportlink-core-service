package mapper

import (
	"sportlink/api/domain/match"
	"sportlink/api/infrastructure/rest/match/response"
)

func EntityToResponse(entity match.Entity) response.MatchResponse {
	return response.MatchResponse{
		ID:           entity.ID,
		Participants: entity.Participants,
		Sport:        string(entity.Sport),
		Day:          entity.Day,
		Status:       entity.Status.String(),
		CreatedAt:    entity.CreatedAt,
	}
}
