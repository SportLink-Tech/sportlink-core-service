package match

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/match"
	"strings"
	"time"
)

// Canonical record keys
const canonicalEntityID = "Entity#Match"

func canonicalIDKey(matchID string) string {
	return matchID
}

// Pointer record keys — one per participant account, immutable
func matchAccountEntityID(accountID string) string {
	return "Entity#MatchAccount#" + accountID
}

func matchAccountIDKey(matchID string) string {
	return "Match#" + matchID
}

// MatchDto is the canonical record. It is the single source of truth for all
// mutable match data (status, result, etc.). There is exactly one per match.
type MatchDto struct {
	EntityId         string `dynamodbav:"EntityId"` // "Entity#Match"
	Id               string `dynamodbav:"Id"`       // "<ulid>"
	LocalAccountId   string `dynamodbav:"LocalAccountId"`
	VisitorAccountId string `dynamodbav:"VisitorAccountId"`
	Sport            string `dynamodbav:"Sport"`
	Day              int64  `dynamodbav:"Day"`            // Unix timestamp (start of day UTC)
	Status           string `dynamodbav:"Status"`
	LocalScore       *int   `dynamodbav:"LocalScore"`
	VisitorScore     *int   `dynamodbav:"VisitorScore"`
	WinnerAccountId  string `dynamodbav:"WinnerAccountId"` // empty when not played or draw
	CreatedAt        int64  `dynamodbav:"CreatedAt"`       // Unix timestamp
}

func (d *MatchDto) ToDomain() match.Entity {
	status, _ := match.ParseStatus(d.Status)

	var result *match.Result
	if d.LocalScore != nil && d.VisitorScore != nil {
		result = &match.Result{
			LocalScore:   *d.LocalScore,
			VisitorScore: *d.VisitorScore,
		}
	}

	return match.Entity{
		ID:               strings.TrimPrefix(d.Id, "MatchId#"),
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

// MatchAccountDto is an immutable pointer record, one per participant account.
// It holds no mutable data — its sole purpose is to allow listing a match by account.
type MatchAccountDto struct {
	EntityId string `dynamodbav:"EntityId"` // "Entity#MatchAccount#<accountId>"
	Id       string `dynamodbav:"Id"`       // "Match#<matchId>"
}

// fromEntity builds the canonical MatchDto and the two MatchAccountDto pointers.
func fromEntity(entity match.Entity) (canonical MatchDto, local MatchAccountDto, visitor MatchAccountDto) {
	canonical = MatchDto{
		EntityId:         canonicalEntityID,
		Id:               canonicalIDKey(entity.ID),
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
		canonical.LocalScore = &ls
		canonical.VisitorScore = &vs
	}

	local = MatchAccountDto{
		EntityId: matchAccountEntityID(entity.LocalAccountID),
		Id:       matchAccountIDKey(entity.ID),
	}

	visitor = MatchAccountDto{
		EntityId: matchAccountEntityID(entity.VisitorAccountID),
		Id:       matchAccountIDKey(entity.ID),
	}

	return canonical, local, visitor
}
