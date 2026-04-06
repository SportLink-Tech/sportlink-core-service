package team

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/team/request"
	"sportlink/api/domain/common"
	"sportlink/api/domain/team"

	"github.com/gin-gonic/gin"
)

// UpdateTeam handles PATCH /sport/:sport/team/:team
// Applies a partial update to the team. Only fields present in the body are modified.
func (sc *DefaultController) UpdateTeam(c *gin.Context) {
	sportParam := c.Param("sport")
	teamName := c.Param("team")

	var req request.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	input := team.PatchInput{
		ID: team.ID{
			Sport: common.Sport(sportParam),
			Name:  teamName,
		},
	}
	if req.Name != "" {
		input.Name = &req.Name
	}

	result, err := sc.updateTeamUC.Invoke(c.Request.Context(), input)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}
