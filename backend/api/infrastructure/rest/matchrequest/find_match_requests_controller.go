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
//   - role=requester: returns requests sent by the account (filters by RequesterAccountID)
//   - role=owner (default): returns requests received by the account (filters by OwnerAccountID)
//   - statuses=PENDING,ACCEPTED: optional comma-separated status filter
func (sc *DefaultController) FindMatchRequests(c *gin.Context) {
	accountID := c.Param("account_id")

	query := domain.DomainQuery{}

	if c.Query("role") == "requester" {
		query.RequesterAccountIDs = []string{accountID}
	} else {
		query.OwnerAccountIDs = []string{accountID}
	}

	if raw := c.Query("statuses"); raw != "" {
		statuses, err := parseStatuses(raw)
		if err != nil {
			c.Error(errors.RequestValidationFailed("invalid status: " + raw))
			return
		}
		query.Statuses = statuses
	}

	entities, err := sc.findMatchRequestsUC.Invoke(c.Request.Context(), query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(entities))
}

func parseStatuses(raw string) ([]domain.Status, error) {
	var statuses []domain.Status
	for _, s := range strings.Split(raw, ",") {
		status, err := domain.ParseStatus(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}
