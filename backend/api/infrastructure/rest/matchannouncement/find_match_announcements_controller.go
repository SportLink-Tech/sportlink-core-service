package matchannouncement

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/matchannouncement"
	restmapper "sportlink/api/infrastructure/rest/matchannouncement/mapper"

	"github.com/gin-gonic/gin"
)

// FindMatchAnnouncements handles the GET request to find match announcements.
func (sc *DefaultController) FindMatchAnnouncements(c *gin.Context) {
	query, err := sc.buildDomainQuery(c)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.findMatchAnnouncementsUC.Invoke(c.Request.Context(), query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	// Check if no results
	if result == nil || len(result.Entities) == 0 {
		c.Error(errors.NotFound("No match announcements found matching the search criteria"))
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

// buildDomainQuery builds a DomainQuery from HTTP query parameters using the parser
func (sc *DefaultController) buildDomainQuery(c *gin.Context) (matchannouncement.DomainQuery, error) {
	query := matchannouncement.DomainQuery{}

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
