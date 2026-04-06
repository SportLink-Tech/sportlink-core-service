package matchrequest

import (
	"net/http"
	"sportlink/api/application/errors"
	apprequest "sportlink/api/application/matchrequest/request"
	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/matchrequest"

	"github.com/gin-gonic/gin"
)

// UpdateMatchRequestStatus handles PATCH /account/:account_id/match-request/:request_id
func (sc *DefaultController) UpdateMatchRequestStatus(c *gin.Context) {
	var req apprequest.UpdateMatchRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(req); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	newStatus, err := matchrequest.ParseStatus(req.Status)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	ownerAccountID := c.Param("account_id")
	requestID := c.Param("request_id")

	err = sc.updateMatchRequestStatusUC.Invoke(c.Request.Context(), usecases.UpdateMatchRequestStatusInput{
		ID:             requestID,
		OwnerAccountID: ownerAccountID,
		NewStatus:      newStatus,
	})
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
