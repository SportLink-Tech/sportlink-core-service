package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"

	"github.com/gin-gonic/gin"
)

// DeleteMatchOffer handles DELETE /account/:account_id/match-offer/:offer_id
func (sc *DefaultController) DeleteMatchOffer(c *gin.Context) {
	offerID := c.Param("offer_id")
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
