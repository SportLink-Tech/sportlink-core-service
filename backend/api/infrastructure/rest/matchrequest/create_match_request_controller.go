package matchrequest

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/matchrequest/usecases"
	restmapper "sportlink/api/infrastructure/rest/matchrequest/mapper"

	"github.com/gin-gonic/gin"
)

// CreateMatchRequest handles POST /account/:account_id/match-offer/:offer_id/match-request
func (sc *DefaultController) CreateMatchRequest(c *gin.Context) {
	requesterAccountID := c.Param("account_id")
	matchOfferID := c.Param("offer_id")

	result, err := sc.createMatchRequestUC.Invoke(c.Request.Context(), usecases.CreateMatchRequestInput{
		MatchOfferID: matchOfferID,
		RequesterAccountID:  requesterAccountID,
	})
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, restmapper.EntityToResponse(*result))
}
