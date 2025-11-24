package rest

import (
	"context"
	"log"
	uplayer "sportlink/api/application/player/usecases"
	uteam "sportlink/api/application/team/usecases"
	"sportlink/api/infrastructure/config"
	iplayer "sportlink/api/infrastructure/persistence/player"
	iteam "sportlink/api/infrastructure/persistence/team"

	"github.com/gin-gonic/gin"

	"sportlink/api/infrastructure/rest/monitoring"
	cplayer "sportlink/api/infrastructure/rest/player"
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

	// Repositories
	playerRepository := iplayer.NewDynamoDBRepository(dynamoDbClient, "SportLinkCore")
	teamRepository := iteam.NewRepository(dynamoDbClient, "SportLinkCore")

	// Player Use Cases
	createPlayer := uplayer.NewCreatePlayerUC(playerRepository)

	// Team Use Cases
	createTeam := uteam.NewCreateTeamUC(playerRepository, teamRepository)
	retrieveTeam := uteam.NewRetrieveTeamUC(teamRepository)
	findTeam := uteam.NewFindTeamUC(teamRepository)

	// Controllers
	playerController := cplayer.NewController(&createPlayer, customValidator)
	router.POST("/player", playerController.CreatePlayer)

	teamController := cteam.NewController(createTeam, retrieveTeam, findTeam, customValidator)
	router.POST("/team", teamController.CreateTeam)
	router.GET("/sport/:sport/team/:team", teamController.RetrieveTeam)
	router.GET("/sport/:sport/team", teamController.FindTeam)

	monitoring.RegisterMetricsRoute(router)
}
