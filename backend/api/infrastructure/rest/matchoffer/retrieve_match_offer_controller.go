package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/matchoffer"
	restmapper "sportlink/api/infrastructure/rest/matchoffer/mapper"

	"github.com/gin-gonic/gin"
)

// RetrieveMatchOffer handles GET /match-offer/:offer_id
func (sc *DefaultController) RetrieveMatchOffer(c *gin.Context) {
	offerID := c.Param("offer_id")

	entity, err := sc.retrieveMatchOfferUC.Invoke(c.Request.Context(), matchoffer.DomainQuery{
		IDs: []string{offerID},
	})
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}
	if entity == nil {
		c.Error(errors.NotFound("match offer not found"))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntityToResponse(*entity))
}
