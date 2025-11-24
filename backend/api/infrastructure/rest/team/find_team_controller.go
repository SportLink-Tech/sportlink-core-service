package team

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/common"
	"sportlink/api/domain/team"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// FindTeam handles the GET request to find teams by sport, name pattern, and categories.
// Endpoint: GET /sport/:sport/team?name=<name>&category=<category1,category2>
// All query parameters are optional (sport in path is required)
func (sc *DefaultController) FindTeam(c *gin.Context) {
	sportParam := c.Param("sport")
	nameQuery := c.Query("name")
	categoryQuery := c.Query("category")

	if sportParam == "" {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	query := team.DomainQuery{
		Sports: []common.Sport{common.Sport(sportParam)},
		Name:   nameQuery,
	}

	// Parse categories if provided
	if categoryQuery != "" {
		categoryStrings := strings.Split(categoryQuery, ",")
		var categories []common.Category
		for _, catStr := range categoryStrings {
			catInt, err := strconv.Atoi(strings.TrimSpace(catStr))
			if err != nil {
				c.Error(errors.RequestValidationFailed("invalid category format"))
				return
			}
			category, err := common.GetCategory(catInt)
			if err != nil {
				c.Error(errors.RequestValidationFailed(err.Error()))
				return
			}
			categories = append(categories, category)
		}
		query.Categories = categories
	}

	teams, err := sc.findTeamUC.Invoke(query)

	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	// Verificar si no se encontraron equipos
	if teams == nil || len(*teams) == 0 {
		c.Error(errors.NotFound("No teams found matching the search criteria"))
		return
	}

	c.JSON(http.StatusOK, *teams)
}
