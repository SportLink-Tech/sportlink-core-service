package matchrequest

import (
	"net/http"
	"sportlink/api/application/errors"
	restmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"

	"github.com/gin-gonic/gin"
)

// FindMatchRequests handles GET /account/:accountId/match-request
func (sc *DefaultController) FindMatchRequests(c *gin.Context) {
	ownerAccountID := c.Param("accountId")

	entities, err := sc.findMatchRequestsUC.Invoke(c.Request.Context(), ownerAccountID)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(entities))
}
