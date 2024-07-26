package rest

import (
	"github.com/gin-gonic/gin"
	uteam "sportlink/api/app/usecase/team"
	"sportlink/api/app/validator"
	cteam "sportlink/api/infrastructure/rest/team"
)

func Routes(router *gin.Engine) {
	router.GET("/livez", livenessHandler)
	router.GET("/readyz", readinessHandler)

	customValidator := validator.GetInstance()

	// Team
	createTeam := uteam.NewCreateTeamUC()
	teamController := cteam.NewController(createTeam, customValidator)

	router.POST("/team", teamController.TeamCreationHandler)
	RegisterMetricsRoute(router)
}
