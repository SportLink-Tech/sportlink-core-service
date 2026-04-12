package matchrequest

import (
	"net/http"
	"sportlink/api/application/matchrequest/usecases"
	reqmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"

	"github.com/gin-gonic/gin"
)

// AcceptMatchRequest handles POST /account/:account_id/match-request/:request_id/accept
func (sc *DefaultController) AcceptMatchRequest(c *gin.Context) {
	ownerAccountID := c.Param("account_id")
	requestID := c.Param("request_id")

	result, err := sc.acceptMatchRequestUC.Invoke(c.Request.Context(), usecases.AcceptMatchRequestInput{
		MatchRequestId: requestID,
		OwnerAccountID: ownerAccountID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reqmapper.EntityToResponse(*result))
}
