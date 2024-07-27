package team

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sportlink/api/app/errors"
	"sportlink/api/app/usecase"
	"sportlink/api/domain/team"
)

// NewTeamRequest defines the structure of the request body for the team creation endpoint.
type NewTeamRequest struct {
	Sport    string `json:"sport" validate:"required,oneof=football paddle"`
	Name     string `json:"name" validate:"required"`
	Category int    `json:"category" validate:"omitempty,category"`
}

type Controller struct {
	createTeamUc usecase.UseCase[NewTeamRequest, team.Entity]
	validator    *validator.Validate
}

func NewController(
	createTeamUc usecase.UseCase[NewTeamRequest, team.Entity],
	validator *validator.Validate,
) *Controller {
	return &Controller{
		createTeamUc: createTeamUc,
		validator:    validator,
	}
}

// TeamCreationHandler handles the POST request to add or modify sports data.
func (sc *Controller) TeamCreationHandler(c *gin.Context) {
	var request NewTeamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(request); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.createTeamUc.Invoke(request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}