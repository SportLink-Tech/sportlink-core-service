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

	announcements, err := sc.findMatchAnnouncementsUC.Invoke(query)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	// Check if no results
	if announcements == nil || len(*announcements) == 0 {
		c.Error(errors.NotFound("No match announcements found matching the search criteria"))
		return
	}

	// Convert domain entities to response DTOs
	responseDTOs := restmapper.EntitiesToResponses(*announcements)

	c.JSON(http.StatusOK, responseDTOs)
}

// buildDomainQuery builds a DomainQuery from HTTP query parameters using the parser
func (sc *DefaultController) buildDomainQuery(c *gin.Context) (matchannouncement.DomainQuery, error) {
	query := matchannouncement.DomainQuery{}

	// Parse sports
	sports, err := sc.queryParser.ParseSports(c.Query("sports"))
	if err != nil {
		return query, err
	}
	query.Sports = sports

	// Parse categories
	categories, err := sc.queryParser.ParseCategories(c.Query("categories"))
	if err != nil {
		return query, err
	}
	query.Categories = categories

	// Parse statuses
	statuses, err := sc.queryParser.ParseStatuses(c.Query("statuses"))
	if err != nil {
		return query, err
	}
	query.Statuses = statuses

	// Parse fromDate
	fromDate, err := sc.queryParser.ParseDate(c.Query("fromDate"))
	if err != nil {
		return query, err
	}
	if !fromDate.IsZero() {
		query.FromDate = fromDate
	}

	// Parse toDate
	toDate, err := sc.queryParser.ParseDate(c.Query("toDate"))
	if err != nil {
		return query, err
	}
	if !toDate.IsZero() {
		query.ToDate = toDate
	}

	// Parse location
	query.Location = sc.queryParser.ParseLocation(
		c.Query("country"),
		c.Query("province"),
		c.Query("locality"),
	)

	return query, nil
}
