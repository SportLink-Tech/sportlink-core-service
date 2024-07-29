package rest

import (
	"github.com/gin-gonic/gin"
	uteam "sportlink/api/application/team/usecases"
	"sportlink/api/infrastructure/config"
	iplayer "sportlink/api/infrastructure/persistence/player"
	"sportlink/api/infrastructure/rest/monitoring"
	cteam "sportlink/api/infrastructure/rest/team"
	"sportlink/api/infrastructure/validator"
)

func Routes(router *gin.Engine) {
	router.GET("/livez", monitoring.LivenessHandler)
	router.GET("/readyz", monitoring.ReadinessHandler)

	customValidator := validator.GetInstance()
	dynamoDbClient := config.NewDynamoDBClient()

	// Player
	playerRepository := iplayer.NewDynamoDBRepository(dynamoDbClient, "SportLinkCore")

	// Team
	createTeam := uteam.NewCreateTeamUC(playerRepository)
	teamController := cteam.NewController(createTeam, customValidator)

	router.POST("/team", teamController.TeamCreationHandler)

	monitoring.RegisterMetricsRoute(router)
}
