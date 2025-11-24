package player

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/player/mapper"
	"sportlink/api/application/player/request"
)

// CreatePlayer handles the POST request to create a new player.
func (sc *DefaultController) CreatePlayer(c *gin.Context) {
	var newPlayerRequest request.NewPlayerRequest
	if err := c.ShouldBindJSON(&newPlayerRequest); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(newPlayerRequest); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	input, err := mapper.CreationRequestToEntity(newPlayerRequest)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.createPlayerUC.Invoke(input)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, result)
}
