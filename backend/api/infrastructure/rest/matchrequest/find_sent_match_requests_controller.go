package matchrequest

import (
	"net/http"
	"sportlink/api/application/errors"
	domain "sportlink/api/domain/matchrequest"
	restmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"
	"strings"

	"github.com/gin-gonic/gin"
)

// FindSentMatchRequests handles GET /account/:accountId/sent-match-request
// Returns match requests sent BY the given account, optionally filtered by status.
func (sc *DefaultController) FindSentMatchRequests(c *gin.Context) {
	requesterAccountID := c.Param("accountId")

	var statuses []domain.Status
	if raw := c.Query("statuses"); raw != "" {
		for _, s := range strings.Split(raw, ",") {
			status, err := domain.ParseStatus(strings.TrimSpace(s))
			if err != nil {
				c.Error(errors.RequestValidationFailed("invalid status: " + s))
				return
			}
			statuses = append(statuses, status)
		}
	}

	entities, err := sc.findSentMatchRequestsUC.Invoke(c.Request.Context(), requesterAccountID, statuses)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(entities))
}
