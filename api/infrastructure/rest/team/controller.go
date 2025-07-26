package team

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"sportlink/api/application"
	"sportlink/api/domain/team"
)

type Controller interface {
	CreateTeam(c *gin.Context)
	RetrieveTeam(c *gin.Context)
}

type DefaultController struct {
	createTeamUC   application.UseCase[team.Entity, team.Entity]
	retrieveTeamUC application.UseCase[team.ID, team.Entity]
	validator      *validator.Validate
}

func NewController(
	createTeamUc application.UseCase[team.Entity, team.Entity],
	retrieveTeamUC application.UseCase[team.ID, team.Entity],
	validator *validator.Validate,
) Controller {
	return &DefaultController{
		createTeamUC:   createTeamUc,
		retrieveTeamUC: retrieveTeamUC,
		validator:      validator,
	}
}
