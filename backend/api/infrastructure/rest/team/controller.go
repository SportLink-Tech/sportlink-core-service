package team

import (
	"sportlink/api/application"
	"sportlink/api/domain/team"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	CreateTeam(c *gin.Context)
	RetrieveTeam(c *gin.Context)
	FindTeam(c *gin.Context)
}

type DefaultController struct {
	createTeamUC   application.UseCase[team.Entity, team.Entity]
	retrieveTeamUC application.UseCase[team.ID, team.Entity]
	findTeamUC     application.UseCase[team.DomainQuery, []team.Entity]
	validator      *validator.Validate
}

func NewController(
	createTeamUc application.UseCase[team.Entity, team.Entity],
	retrieveTeamUC application.UseCase[team.ID, team.Entity],
	findTeamUC application.UseCase[team.DomainQuery, []team.Entity],
	validator *validator.Validate,
) Controller {
	return &DefaultController{
		createTeamUC:   createTeamUc,
		retrieveTeamUC: retrieveTeamUC,
		findTeamUC:     findTeamUC,
		validator:      validator,
	}
}
