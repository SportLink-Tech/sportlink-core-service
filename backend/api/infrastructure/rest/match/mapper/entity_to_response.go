package mapper

import (
	"sportlink/api/domain/match"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/infrastructure/rest/match/response"
)

func EntityToResponse(entity match.Entity, offer *matchoffer.Entity) response.MatchResponse {
	r := response.MatchResponse{
		ID:           entity.ID,
		Participants: entity.Participants,
		Sport:        string(entity.Sport),
		Day:          entity.Day,
		Status:       entity.Status.String(),
		CreatedAt:    entity.CreatedAt,
	}

	if offer != nil {
		r.Title = offer.GetTitle()
		r.TimeSlot = &response.TimeSlotResponse{
			StartTime: offer.TimeSlot.StartTime,
			EndTime:   offer.TimeSlot.EndTime,
		}
	}

	return r
}
