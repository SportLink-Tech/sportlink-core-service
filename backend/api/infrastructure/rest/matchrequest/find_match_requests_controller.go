package matchrequest

import (
	"net/http"
	"sportlink/api/application/errors"
	domain "sportlink/api/domain/matchrequest"
	restmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"
	"strings"

	"github.com/gin-gonic/gin"
)

// FindMatchRequests handles GET /account/:accountId/match-request
// When ?sent=true is provided, returns match requests sent BY the account.
// Otherwise, returns match requests received by the account.
func (sc *DefaultController) FindMatchRequests(c *gin.Context) {
	accountID := c.Param("accountId")

	if c.Query("sent") == "true" {
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

		entities, err := sc.findSentMatchRequestsUC.Invoke(c.Request.Context(), accountID, statuses)
		if err != nil {
			c.Error(errors.UseCaseExecutionFailed(err.Error()))
			return
		}

		c.JSON(http.StatusOK, restmapper.EntitiesToResponses(entities))
		return
	}

	entities, err := sc.findMatchRequestsUC.Invoke(c.Request.Context(), accountID)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(entities))
}
