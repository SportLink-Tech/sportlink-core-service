package mapper

import (
	"sportlink/api/application/player/request"
	"sportlink/api/domain/common"
	"sportlink/api/domain/player"
)

func CreationRequestToEntity(req request.NewPlayerRequest) (player.Entity, error) {
	category, err := common.GetCategory(req.Category)
	if err != nil {
		return player.Entity{}, err
	}

	sport := common.Sport(req.Sport)

	return player.NewPlayer(category, sport), nil
}
