package match

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/match"
	"time"
)

// entityIDPrefix builds the DynamoDB partition key for a given account.
// Each match is stored under two partition keys — one per participant — so
// both accounts can query their matches without a scan.
func entityIDPrefix(accountID string) string {
	return "AccountMatch#" + accountID
}

type Dto struct {
	EntityId         string  `dynamodbav:"EntityId"`         // "AccountMatch#<account_id>"
	Id               string  `dynamodbav:"Id"`               // ULID — same for both records of the same match
	LocalAccountId   string  `dynamodbav:"LocalAccountId"`
	VisitorAccountId string  `dynamodbav:"VisitorAccountId"`
	Sport            string  `dynamodbav:"Sport"`
	Day              int64   `dynamodbav:"Day"`              // Unix timestamp (start of day UTC)
	Status           string  `dynamodbav:"Status"`
	LocalScore       *int    `dynamodbav:"LocalScore"`
	VisitorScore     *int    `dynamodbav:"VisitorScore"`
	WinnerAccountId  string  `dynamodbav:"WinnerAccountId"`  // empty when not played or draw
	CreatedAt        int64   `dynamodbav:"CreatedAt"`        // Unix timestamp
}

func (d *Dto) ToDomain() match.Entity {
	status, _ := match.ParseStatus(d.Status)

	var result *match.Result
	if d.LocalScore != nil && d.VisitorScore != nil {
		result = &match.Result{
			LocalScore:   *d.LocalScore,
			VisitorScore: *d.VisitorScore,
		}
	}

	return match.Entity{
		ID:               d.Id,
		LocalAccountID:   d.LocalAccountId,
		VisitorAccountID: d.VisitorAccountId,
		Sport:            common.Sport(d.Sport),
		Day:              time.Unix(d.Day, 0).UTC(),
		Status:           status,
		Result:           result,
		WinnerAccountID:  d.WinnerAccountId,
		CreatedAt:        time.Unix(d.CreatedAt, 0).UTC(),
	}
}

// fromEntity builds the two DTOs (one per account) that represent a match in DynamoDB.
func fromEntity(entity match.Entity) (local Dto, visitor Dto) {
	base := Dto{
		Id:               entity.ID,
		LocalAccountId:   entity.LocalAccountID,
		VisitorAccountId: entity.VisitorAccountID,
		Sport:            string(entity.Sport),
		Day:              entity.Day.Unix(),
		Status:           entity.Status.String(),
		WinnerAccountId:  entity.WinnerAccountID,
		CreatedAt:        entity.CreatedAt.Unix(),
	}

	if entity.Result != nil {
		ls := entity.Result.LocalScore
		vs := entity.Result.VisitorScore
		base.LocalScore = &ls
		base.VisitorScore = &vs
	}

	local = base
	local.EntityId = entityIDPrefix(entity.LocalAccountID)

	visitor = base
	visitor.EntityId = entityIDPrefix(entity.VisitorAccountID)

	return local, visitor
}
