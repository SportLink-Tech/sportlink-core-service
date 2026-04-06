package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/matchoffer"
	restmapper "sportlink/api/infrastructure/rest/matchoffer/mapper"

	"github.com/gin-gonic/gin"
)

// FindMatchOffers handles the GET request to find match offers.
func (sc *DefaultController) FindMatchOffers(c *gin.Context) {
	query, err := sc.buildDomainQuery(c)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.findMatchOffersUC.Invoke(c.Request.Context(), query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	// Convert domain entities to response DTOs
	responseDTOs := restmapper.EntitiesToResponses(result.Entities)

	// Build paginated response
	paginatedResponse := restmapper.NewPaginatedResponse(
		responseDTOs,
		result.Page.Number,
		result.Page.OutOf,
		result.Page.Total,
	)

	c.JSON(http.StatusOK, paginatedResponse)
}

// FindAccountMatchOffers handles GET /account/:accountId/match-offer
// Returns all match offers owned by the given account, optionally filtered by status.
func (sc *DefaultController) FindAccountMatchOffers(c *gin.Context) {
	accountID := c.Param("accountId")

	statuses, err := sc.queryParser.Statuses(c.Query("statuses"))
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	query := matchoffer.DomainQuery{
		OwnerAccountID: accountID,
		Statuses:       statuses,
	}

	result, err := sc.findMatchOffersUC.Invoke(c.Request.Context(), query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, restmapper.EntitiesToResponses(result.Entities))
}

// buildDomainQuery builds a DomainQuery from HTTP query parameters using the parser
func (sc *DefaultController) buildDomainQuery(c *gin.Context) (matchoffer.DomainQuery, error) {
	query := matchoffer.DomainQuery{}

	// Parse sports
	sports, err := sc.queryParser.Sports(c.Query("sports"))
	if err != nil {
		return query, err
	}
	query.Sports = sports

	// Parse categories
	categories, err := sc.queryParser.Categories(c.Query("categories"))
	if err != nil {
		return query, err
	}
	query.Categories = categories

	// Parse statuses
	statuses, err := sc.queryParser.Statuses(c.Query("statuses"))
	if err != nil {
		return query, err
	}
	query.Statuses = statuses

	// Parse fromDate
	fromDate, err := sc.queryParser.Date(c.Query("fromDate"))
	if err != nil {
		return query, err
	}
	if !fromDate.IsZero() {
		query.FromDate = fromDate
	}

	// Parse toDate
	toDate, err := sc.queryParser.Date(c.Query("toDate"))
	if err != nil {
		return query, err
	}
	if !toDate.IsZero() {
		query.ToDate = toDate
	}

	// Parse location
	query.Location = sc.queryParser.Location(
		c.Query("country"),
		c.Query("province"),
		c.Query("locality"),
	)

	// Parse geo filter (takes precedence over text location when present)
	geoFilter, err := sc.queryParser.GeoFilter(
		c.Query("lat"),
		c.Query("lng"),
		c.Query("radiusKm"),
	)
	if err != nil {
		return query, err
	}
	if geoFilter != nil {
		query.GeoFilter = geoFilter
		query.Location = nil // geo filter replaces text location filter
	}

	// Parse limit
	limit, err := sc.queryParser.Limit(c.Query("limit"))
	if err != nil {
		return query, err
	}
	query.Limit = limit

	// Parse offset
	offset, err := sc.queryParser.Offset(c.Query("offset"))
	if err != nil {
		return query, err
	}
	query.Offset = offset

	return query, nil
}
