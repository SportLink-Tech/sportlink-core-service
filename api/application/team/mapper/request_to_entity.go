package mapper

import (
	team2 "sportlink/api/application/team/request"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
)

func CreationRequestToEntity(request team2.NewTeamRequest) (team.Entity, error) {
	category, err := common.GetCategory(request.Category)
	if err != nil {
		return team.Entity{}, err
	}
	stats := common.NewStats(0, 0, 0)
	players := make([]player.Entity, 0)
	sport := common.Sport(request.Sport)
	return *team.NewTeam(request.Name, category, *stats, sport, players), nil
}