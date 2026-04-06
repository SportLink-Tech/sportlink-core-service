package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"

	"github.com/gin-gonic/gin"
)

// DeleteMatchOffer handles DELETE /account/:accountId/match-offer/:offerId
func (sc *DefaultController) DeleteMatchOffer(c *gin.Context) {
	offerID := c.Param("offerId")
	if offerID == "" {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	err := sc.deleteMatchOfferUC.Invoke(c.Request.Context(), offerID)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
