package team

import (
	"github.com/go-playground/validator/v10"
	"sportlink/api/application"
	"sportlink/api/domain/team"
)

type Controller struct {
	createTeamUC application.UseCase[team.Entity, team.Entity]
	validator    *validator.Validate
}

func NewController(
	createTeamUc application.UseCase[team.Entity, team.Entity],
	validator *validator.Validate,
) *Controller {
	return &Controller{
		createTeamUC: createTeamUc,
		validator:    validator,
	}
}
