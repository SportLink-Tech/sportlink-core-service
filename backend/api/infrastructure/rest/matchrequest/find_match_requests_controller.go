package matchrequest

import (
	"net/http"
	"sportlink/api/application/errors"
	domain "sportlink/api/domain/matchrequest"
	restmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"
	"strings"

	"github.com/gin-gonic/gin"
)

// FindMatchRequests handles GET /account/:account_id/match-request
// Query params:
//   - sent=true: returns requests sent BY the account (filters by RequesterAccountID)
//   - statuses=pending,accepted: optional status filter (only applied when sent=true)
//
// Default: returns requests received by the account (filters by OwnerAccountID)
func (sc *DefaultController) FindMatchRequests(c *gin.Context) {
	accountID := c.Param("account_id")

	query := domain.DomainQuery{}

	if c.Query("sent") == "true" {
		query.RequesterAccountIDs = []string{accountID}

		if raw := c.Query("statuses"); raw != "" {
			for _, s := range strings.Split(raw, ",") {
				status, err := domain.ParseStatus(strings.TrimSpace(s))
				if err != nil {
					c.Error(errors.RequestValidationFailed("invalid status: " + s))
					return
				}
				query.Statuses = append(query.Statuses, status)
			}
		}
	} else {
		query.OwnerAccountIDs = []string{accountID}
	}

	entities, err := sc.findMatchRequestsUC.Invoke(c.Request.Context(), query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(entities))
}
