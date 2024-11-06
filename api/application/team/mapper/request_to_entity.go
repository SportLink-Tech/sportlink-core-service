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
	sport := common.Sport(request.Sport)
	players := make([]player.Entity, 0)
	for _, playerId := range request.PlayerIds {
		players = append(players, player.Entity{ID: playerId})
	}
	return *team.NewTeam(request.Name, category, *stats, sport, players), nil
}
