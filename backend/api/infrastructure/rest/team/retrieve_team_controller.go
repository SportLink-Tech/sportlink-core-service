package team

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/domain/common"
	"sportlink/api/domain/team"
)

// RetrieveTeam handles the GET request to retrieve a team by sport and name.
// Endpoint: GET /sport/:sport/team/:team
func (sc *DefaultController) RetrieveTeam(c *gin.Context) {
	sportParam := c.Param("sport")
	teamName := c.Param("team")

	if sportParam == "" || teamName == "" {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	t, err := sc.retrieveTeamUC.Invoke(c.Request.Context(), team.ID{
		Sport: common.Sport(sportParam),
		Name:  teamName,
	})

	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, t)
}
