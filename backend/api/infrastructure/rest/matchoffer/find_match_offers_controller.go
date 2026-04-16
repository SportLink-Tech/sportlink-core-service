package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/matchoffer"
	restmapper "sportlink/api/infrastructure/rest/matchoffer/mapper"

	"github.com/gin-gonic/gin"
)

// FindAccountMatchOffers handles GET /account/:account_id/match-offer
// Returns all match offers owned by the given account, optionally filtered by status.
func (sc *DefaultController) FindAccountMatchOffers(c *gin.Context) {
	accountID := c.Param("account_id")
	statuses, err := sc.queryParser.Statuses(c.Query("statuses"))
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.findAccountMatchOffersUC.Invoke(c.Request.Context(), matchoffer.DomainQuery{
		OwnerAccountID: accountID,
		Statuses:       statuses,
	})
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(result.Entities))
}

// buildDomainQuery builds a DomainQuery from HTTP query parameters using the parser.
func (sc *DefaultController) buildDomainQuery(c *gin.Context) (matchoffer.DomainQuery, error) {
	query := matchoffer.DomainQuery{}

	sports, err := sc.queryParser.Sports(c.Query("sports"))
	if err != nil {
		return query, err
	}
	query.Sports = sports

	categories, err := sc.queryParser.Categories(c.Query("categories"))
	if err != nil {
		return query, err
	}
	query.Categories = categories

	statuses, err := sc.queryParser.Statuses(c.Query("statuses"))
	if err != nil {
		return query, err
	}
	query.Statuses = statuses

	fromDate, err := sc.queryParser.Date(c.Query("from_date"))
	if err != nil {
		return query, err
	}
	if !fromDate.IsZero() {
		query.FromDate = fromDate
	}

	toDate, err := sc.queryParser.Date(c.Query("to_date"))
	if err != nil {
		return query, err
	}
	if !toDate.IsZero() {
		query.ToDate = toDate
	}

	query.Location = sc.queryParser.Location(
		c.Query("country"),
		c.Query("province"),
		c.Query("locality"),
	)

	geoFilter, err := sc.queryParser.GeoFilter(
		c.Query("lat"),
		c.Query("lng"),
		c.Query("radius_km"),
	)
	if err != nil {
		return query, err
	}
	if geoFilter != nil {
		query.GeoFilter = geoFilter
		query.Location = nil
	}

	limit, err := sc.queryParser.Limit(c.Query("limit"))
	if err != nil {
		return query, err
	}
	query.Limit = limit

	offset, err := sc.queryParser.Offset(c.Query("offset"))
	if err != nil {
		return query, err
	}
	query.Offset = offset

	return query, nil
}
