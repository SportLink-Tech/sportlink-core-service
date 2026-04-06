package rest

import (
	"context"
	"log"
	umatchoffer "sportlink/api/application/matchoffer/usecases"
	umatchrequest "sportlink/api/application/matchrequest/usecases"
	uplayer "sportlink/api/application/player/usecases"
	uteam "sportlink/api/application/team/usecases"
	"sportlink/api/infrastructure/config"
	imatchoffer "sportlink/api/infrastructure/persistence/matchoffer"
	imatchrequest "sportlink/api/infrastructure/persistence/matchrequest"
	iplayer "sportlink/api/infrastructure/persistence/player"
	iteam "sportlink/api/infrastructure/persistence/team"

	"github.com/gin-gonic/gin"

	cmatchoffer "sportlink/api/infrastructure/rest/matchoffer"
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
	matchOfferRepository := imatchoffer.NewRepository(dynamoDbClient, "SportLinkCore")
	matchRequestRepository := imatchrequest.NewRepository(dynamoDbClient, "SportLinkCore")

	// Player Use Cases
	createPlayer := uplayer.NewCreatePlayerUC(playerRepository)

	// Team Use Cases
	createTeam := uteam.NewCreateTeamUC(playerRepository, teamRepository)
	retrieveTeam := uteam.NewRetrieveTeamUC(teamRepository)
	findTeam := uteam.NewFindTeamUC(teamRepository)

	// Match Offer Use Cases
	createMatchOffer := umatchoffer.NewCreateMatchOfferUC(matchOfferRepository, teamRepository)
	findMatchOffers := umatchoffer.NewFindMatchOfferUC(matchOfferRepository)
	deleteMatchOffer := umatchoffer.NewDeleteMatchOfferUC(matchOfferRepository)

	// Match Request Use Cases
	createMatchRequest := umatchrequest.NewCreateMatchRequestUC(matchRequestRepository, matchOfferRepository)
	findMatchRequests := umatchrequest.NewFindMatchRequestsUC(matchRequestRepository)
	findSentMatchRequests := umatchrequest.NewFindSentMatchRequestsUC(matchRequestRepository)
	updateMatchRequestStatus := umatchrequest.NewUpdateMatchRequestStatusUC(matchRequestRepository)

	// Controllers
	playerController := cplayer.NewController(&createPlayer, customValidator)
	router.POST("/player", playerController.CreatePlayer)

	teamController := cteam.NewController(createTeam, retrieveTeam, findTeam, findTeam, customValidator)
	router.POST("/account/:accountId/team", teamController.CreateTeam)
	router.GET("/account/:accountId/team", teamController.ListAccountTeams)
	router.GET("/sport/:sport/team/:team", teamController.RetrieveTeam)
	router.GET("/sport/:sport/team", teamController.FindTeam)

	matchOfferController := cmatchoffer.NewController(createMatchOffer, findMatchOffers, deleteMatchOffer, customValidator)
	router.POST("/account/:accountId/match-offer", matchOfferController.CreateMatchOffer)
	router.GET("/match-offer", matchOfferController.FindMatchOffers)
	router.GET("/account/:accountId/match-offer", matchOfferController.FindAccountMatchOffers)
	router.DELETE("/account/:accountId/match-offer/:offerId", matchOfferController.DeleteMatchOffer)

	matchRequestController := cmatchrequest.NewController(createMatchRequest, findMatchRequests, findSentMatchRequests, updateMatchRequestStatus, customValidator)
	router.POST("/account/:accountId/match-offer/:announcementId/match-request", matchRequestController.CreateMatchRequest)
	router.GET("/account/:accountId/match-request", matchRequestController.FindMatchRequests)
	router.GET("/account/:accountId/sent-match-request", matchRequestController.FindSentMatchRequests)
	router.PATCH("/account/:accountId/match-request/:requestId", matchRequestController.UpdateMatchRequestStatus)

	monitoring.RegisterMetricsRoute(router)
}
