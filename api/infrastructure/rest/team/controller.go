package team

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/app/errors"
	"sportlink/api/app/usecase"
	"sportlink/api/domain/team"
)

// NewTeamRequest defines the structure of the request body for the team creation endpoint.
type NewTeamRequest struct {
	Sport    string `json:"sport" validate:"required,oneof=football paddle"`
	Name     string `json:"name" validate:"required"`
	Category *int   `json:"category" validate:"omitempty,category"`
}

type Controller struct {
	createTeamUc usecase.UseCase[NewTeamRequest, team.Entity]
}

func NewController(
	createTeamUc usecase.UseCase[NewTeamRequest, team.Entity],
) *Controller {
	return &Controller{
		createTeamUc: createTeamUc,
	}
}

// TeamCreationHandler handles the POST request to add or modify sports data.
func (sc *Controller) TeamCreationHandler(c *gin.Context) {
	var request NewTeamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.InvalidRequestData())
		return
	}

	result, err := sc.createTeamUc.Invoke(request)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error.go": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}
