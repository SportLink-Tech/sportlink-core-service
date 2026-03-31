package matchrequest

import (
	"sportlink/api/domain/matchrequest"
	"time"
)

type Dto struct {
	EntityId            string `dynamodbav:"EntityId"`            // "Entity#MatchRequest"
	Id                  string `dynamodbav:"Id"`                  // Generated ULID
	MatchAnnouncementId string `dynamodbav:"MatchAnnouncementId"` // Referenced announcement ID
	OwnerAccountId      string `dynamodbav:"OwnerAccountId"`      // Announcement owner account ID (GSI partition key)
	RequesterAccountId  string `dynamodbav:"RequesterAccountId"`  // Requester account ID
	Status              string `dynamodbav:"Status"`              // PENDING, ACCEPTED, REJECTED
	CreatedAt           int64  `dynamodbav:"CreatedAt"`           // Unix timestamp
}

func (d *Dto) ToDomain() matchrequest.Entity {
	status, _ := matchrequest.ParseStatus(d.Status)

	return matchrequest.Entity{
		ID:                  d.Id,
		MatchAnnouncementID: d.MatchAnnouncementId,
		OwnerAccountID:      d.OwnerAccountId,
		RequesterAccountID:  d.RequesterAccountId,
		Status:              status,
		CreatedAt:           time.Unix(d.CreatedAt, 0).UTC(),
	}
}

func From(entity matchrequest.Entity) Dto {
	return Dto{
		EntityId:            "Entity#MatchRequest",
		Id:                  entity.ID,
		MatchAnnouncementId: entity.MatchAnnouncementID,
		OwnerAccountId:      entity.OwnerAccountID,
		RequesterAccountId:  entity.RequesterAccountID,
		Status:              entity.Status.String(),
		CreatedAt:           entity.CreatedAt.Unix(),
	}
}
