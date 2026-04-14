package matchrequest

import (
	"net/http"
	"sportlink/api/application/matchrequest/usecases"
	reqmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"

	"github.com/gin-gonic/gin"
)

// CancelMatchRequest handles POST /account/:account_id/match-request/:request_id/cancel
func (sc *DefaultController) CancelMatchRequest(c *gin.Context) {
	requesterAccountID := c.Param("account_id")
	requestID := c.Param("request_id")

	result, err := sc.cancelMatchRequestUC.Invoke(c.Request.Context(), usecases.CancelMatchRequestInput{
		MatchRequestId:     requestID,
		RequesterAccountID: requesterAccountID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reqmapper.EntityToResponse(*result))
}
