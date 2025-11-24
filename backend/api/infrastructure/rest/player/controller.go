package player

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"sportlink/api/application"
	"sportlink/api/domain/player"
)

type Controller interface {
	CreatePlayer(c *gin.Context)
}

type DefaultController struct {
	createPlayerUC application.UseCase[player.Entity, player.Entity]
	validator      *validator.Validate
}

func NewController(
	createPlayerUC application.UseCase[player.Entity, player.Entity],
	validator *validator.Validate,
) Controller {
	return &DefaultController{
		createPlayerUC: createPlayerUC,
		validator:      validator,
	}
}
