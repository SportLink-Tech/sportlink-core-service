package team

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/team/mapper"
	"sportlink/api/application/team/request"
)

// CreateTeam handles the POST request to add or modify sports data.
func (sc *DefaultController) CreateTeam(c *gin.Context) {
	var newTeamRequest request.NewTeamRequest
	if err := c.ShouldBindJSON(&newTeamRequest); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(newTeamRequest); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	input, err := mapper.CreationRequestToEntity(newTeamRequest)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.createTeamUC.Invoke(input)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}
