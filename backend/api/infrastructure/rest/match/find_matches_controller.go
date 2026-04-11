package match

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/match/usecases"
	domainmatch "sportlink/api/domain/match"
	"sportlink/api/infrastructure/rest/match/mapper"
	"sportlink/api/infrastructure/rest/match/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// FindMatches handles GET /account/:account_id/match
// Query params:
//   - statuses=ACCEPTED,PLAYED,CANCELLED: optional status filter
func (sc *DefaultController) FindMatches(c *gin.Context) {
	accountID := c.Param("account_id")

	input := usecases.FindMatchesInput{AccountID: accountID}

	if raw := c.Query("statuses"); raw != "" {
		for _, s := range strings.Split(raw, ",") {
			status, err := domainmatch.ParseStatus(strings.TrimSpace(s))
			if err != nil {
				c.Error(errors.RequestValidationFailed("invalid status: " + s))
				return
			}
			input.Statuses = append(input.Statuses, status)
		}
	}

	result, err := sc.findMatchesUC.Invoke(c.Request.Context(), input)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	responses := make([]response.MatchResponse, len(*result))
	for i, e := range *result {
		responses[i] = mapper.EntityToResponse(e)
	}

	c.JSON(http.StatusOK, responses)
}
