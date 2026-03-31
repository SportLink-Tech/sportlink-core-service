package rest

import (
	"context"
	"log"
	umatchannouncement "sportlink/api/application/matchannouncement/usecases"
	umatchrequest "sportlink/api/application/matchrequest/usecases"
	uplayer "sportlink/api/application/player/usecases"
	uteam "sportlink/api/application/team/usecases"
	"sportlink/api/infrastructure/config"
	imatchannouncement "sportlink/api/infrastructure/persistence/matchannouncement"
	imatchrequest "sportlink/api/infrastructure/persistence/matchrequest"
	iplayer "sportlink/api/infrastructure/persistence/player"
	iteam "sportlink/api/infrastructure/persistence/team"

	"github.com/gin-gonic/gin"

	cmatchannouncement "sportlink/api/infrastructure/rest/matchannouncement"
	cmatchrequest "sportlink/api/infrastructure/rest/matchrequest"
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
	matchAnnouncementRepository := imatchannouncement.NewRepository(dynamoDbClient, "SportLinkCore")
	matchRequestRepository := imatchrequest.NewRepository(dynamoDbClient, "SportLinkCore")

	// Player Use Cases
	createPlayer := uplayer.NewCreatePlayerUC(playerRepository)

	// Team Use Cases
	createTeam := uteam.NewCreateTeamUC(playerRepository, teamRepository)
	retrieveTeam := uteam.NewRetrieveTeamUC(teamRepository)
	findTeam := uteam.NewFindTeamUC(teamRepository)

	// Match Announcement Use Cases
	createMatchAnnouncement := umatchannouncement.NewCreateMatchAnnouncementUC(matchAnnouncementRepository, teamRepository)
	findMatchAnnouncements := umatchannouncement.NewFindMatchAnnouncementUC(matchAnnouncementRepository)

	// Match Request Use Cases
	createMatchRequest := umatchrequest.NewCreateMatchRequestUC(matchRequestRepository, matchAnnouncementRepository)
	findMatchRequests := umatchrequest.NewFindMatchRequestsUC(matchRequestRepository)
	updateMatchRequestStatus := umatchrequest.NewUpdateMatchRequestStatusUC(matchRequestRepository)

	// Controllers
	playerController := cplayer.NewController(&createPlayer, customValidator)
	router.POST("/player", playerController.CreatePlayer)

	teamController := cteam.NewController(createTeam, retrieveTeam, findTeam, findTeam, customValidator)
	router.POST("/account/:accountId/team", teamController.CreateTeam)
	router.GET("/account/:accountId/team", teamController.ListAccountTeams)
	router.GET("/sport/:sport/team/:team", teamController.RetrieveTeam)
	router.GET("/sport/:sport/team", teamController.FindTeam)

	matchAnnouncementController := cmatchannouncement.NewController(createMatchAnnouncement, findMatchAnnouncements, customValidator)
	router.POST("/account/:accountId/match-announcement", matchAnnouncementController.CreateMatchAnnouncement)
	router.GET("/match-announcement", matchAnnouncementController.FindMatchAnnouncements)

	matchRequestController := cmatchrequest.NewController(createMatchRequest, findMatchRequests, updateMatchRequestStatus, customValidator)
	router.POST("/account/:accountId/match-announcement/:announcementId/match-request", matchRequestController.CreateMatchRequest)
	router.GET("/account/:accountId/match-request", matchRequestController.FindMatchRequests)
	router.PATCH("/account/:accountId/match-request/:requestId", matchRequestController.UpdateMatchRequestStatus)

	monitoring.RegisterMetricsRoute(router)
}
