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
	retrieveTeam := uteam.NewRetrieveTeamUC(teamRepository)

	// DefaultController
	teamController := cteam.NewController(createTeam, retrieveTeam, customValidator)
	router.POST("/team", teamController.CreateTeam)
	router.GET("/sport/:sport/team/:team", teamController.RetrieveTeam)

	monitoring.RegisterMetricsRoute(router)
}
