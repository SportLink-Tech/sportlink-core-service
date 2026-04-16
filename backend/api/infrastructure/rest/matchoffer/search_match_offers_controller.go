package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/matchoffer/usecases"
	restmapper "sportlink/api/infrastructure/rest/matchoffer/mapper"

	"github.com/gin-gonic/gin"
)

// SearchMatchOffers handles GET /account/:account_id/match-offer/search.
// Returns paginated match offers available for the given account, excluding:
//   - offers owned by the account
//   - offers where the account already has a PENDING or ACCEPTED match request
func (sc *DefaultController) SearchMatchOffers(c *gin.Context) {
	accountID := c.Param("account_id")

	query, err := sc.buildDomainQuery(c)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.searchMatchOffersUC.Invoke(c.Request.Context(), usecases.SearchMatchOffersInput{
		ViewerAccountID: accountID,
		Query:           query,
	})
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	if len(result.Entities) == 0 {
		c.Error(errors.NotFound("no match offers found"))
		return
	}

	responseDTOs := restmapper.EntitiesToResponses(result.Entities)
	paginatedResponse := restmapper.NewPaginatedResponse(
		responseDTOs,
		result.Page.Number,
		result.Page.OutOf,
		result.Page.Total,
	)

	c.JSON(http.StatusOK, paginatedResponse)
}
