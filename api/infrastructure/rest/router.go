package rest

import (
	"github.com/gin-gonic/gin"
	uteam "sportlink/api/application/team/use_cases"
	"sportlink/api/infrastructure/rest/monitoring"
	cteam "sportlink/api/infrastructure/rest/team"
	"sportlink/api/infrastructure/validator"
)

func Routes(router *gin.Engine) {
	router.GET("/livez", monitoring.LivenessHandler)
	router.GET("/readyz", monitoring.ReadinessHandler)

	customValidator := validator.GetInstance()

	// Team
	createTeam := uteam.NewCreateTeamUC()
	teamController := cteam.NewController(createTeam, customValidator)

	router.POST("/team", teamController.TeamCreationHandler)
	monitoring.RegisterMetricsRoute(router)
}
