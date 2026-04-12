package matchoffer

import (
	"net/http"
	"sportlink/api/application/matchoffer/usecases"
	matchmapper "sportlink/api/infrastructure/rest/match/mapper"

	"github.com/gin-gonic/gin"
)

// ConfirmMatchOffer handles POST /account/:account_id/match-offer/:offer_id/confirm
func (c *DefaultController) ConfirmMatchOffer(ctx *gin.Context) {
	ownerAccountID := ctx.Param("account_id")
	offerID := ctx.Param("offer_id")

	result, err := c.confirmMatchOfferUC.Invoke(ctx.Request.Context(), usecases.ConfirmMatchOfferInput{
		MatchOfferID:   offerID,
		OwnerAccountID: ownerAccountID,
	})
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, matchmapper.EntityToResponse(*result))
}
