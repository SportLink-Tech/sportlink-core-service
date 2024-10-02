package team

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/team/mapper"
	request2 "sportlink/api/application/team/request"
)

// TeamCreationHandler handles the POST request to add or modify sports data.
func (sc *Controller) TeamCreationHandler(c *gin.Context) {
	var request request2.NewTeamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(request); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	input, err := mapper.CreationRequestToEntity(request)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.createTeamUC.Invoke(input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, result)
}
