package matchannouncement

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/matchannouncement/mapper"
	"sportlink/api/application/matchannouncement/request"
	restmapper "sportlink/api/infrastructure/rest/matchannouncement/mapper"

	"github.com/gin-gonic/gin"
)

// CreateMatchAnnouncement handles the POST request to create a match announcement.
func (sc *DefaultController) CreateMatchAnnouncement(c *gin.Context) {
	var newMatchAnnouncementRequest request.NewMatchAnnouncementRequest
	if err := c.ShouldBindJSON(&newMatchAnnouncementRequest); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(newMatchAnnouncementRequest); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	input, err := mapper.CreationRequestToEntity(newMatchAnnouncementRequest)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.createMatchAnnouncementUC.Invoke(input)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	// Convert domain entity to response DTO
	responseDTO := restmapper.EntityToResponse(*result)

	c.JSON(http.StatusCreated, responseDTO)
}
