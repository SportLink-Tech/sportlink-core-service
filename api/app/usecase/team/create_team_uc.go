package team

import (
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	dteam "sportlink/api/domain/team"
	"sportlink/api/infrastructure/rest/team"
)

type CreateTeamUC struct {
}

func NewCreateTeamUC() *CreateTeamUC {
	return &CreateTeamUC{}
}

func (uc *CreateTeamUC) Invoke(input team.NewTeamRequest) (*dteam.Entity, error) {
	category := uc.getCategoryOrDefault(input)
	stats := *common.NewStats(0, 0, 0)
	sport := common.Sport(input.Sport)

	return dteam.NewTeam(
		input.Name,
		category,
		stats,
		sport,
		[]player.Entity{},
	), nil
}

func (uc *CreateTeamUC) getCategoryOrDefault(input team.NewTeamRequest) common.Category {
	var category common.Category
	if input.Category == nil {
		category = common.L1
	} else {
		category = common.Category(*input.Category)
	}
	return category
}
