package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	uteam "sportlink/api/application/team/usecases"
	"sportlink/api/infrastructure/config"
	iplayer "sportlink/api/infrastructure/persistence/player"
	iteam "sportlink/api/infrastructure/persistence/team"

	"sportlink/api/infrastructure/rest/monitoring"
	cteam "sportlink/api/infrastructure/rest/team"
	"sportlink/api/infrastructure/validator"
)

func Routes(router *gin.Engine) {
	router.GET("/livez", monitoring.LivenessHandler)
	router.GET("/readyz", monitoring.ReadinessHandler)

	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	customValidator := validator.GetInstance()
	dynamoDbClient := config.NewDynamoDBClient(cfg.DynamoDbCfg)

	// Player
	playerRepository := iplayer.NewDynamoDBRepository(dynamoDbClient, "SportLinkCore")

	// Team
	teamRepository := iteam.NewRepository(dynamoDbClient, "SportLinkCore")

	// Use Cases
	createTeam := uteam.NewCreateTeamUC(playerRepository, teamRepository)

	// Controller
	teamController := cteam.NewController(createTeam, customValidator)

	router.POST("/team", teamController.TeamCreationHandler)

	monitoring.RegisterMetricsRoute(router)
}
