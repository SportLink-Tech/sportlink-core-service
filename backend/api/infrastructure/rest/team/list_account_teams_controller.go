package team

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/team"
)

// ListAccountTeams handles GET /account/:account_id/team
func (sc *DefaultController) ListAccountTeams(c *gin.Context) {
	accountID := c.Param("account_id")
	if accountID == "" {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	query := team.DomainQuery{
		OwnerAccountID: accountID,
	}

	result, err := sc.listAccountTeamsUC.Invoke(c.Request.Context(), query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}
