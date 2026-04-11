package rest

import (
	"context"
	"log"
	uaccount "sportlink/api/application/account/usecases"
	authservice "sportlink/api/application/auth/service"
	uauth "sportlink/api/application/auth/usecases"
	umatchoffer "sportlink/api/application/matchoffer/usecases"
	umatch "sportlink/api/application/match/usecases"
	umatchrequest "sportlink/api/application/matchrequest/usecases"
	imatch "sportlink/api/infrastructure/persistence/match"
	cmatch "sportlink/api/infrastructure/rest/match"
	uplayer "sportlink/api/application/player/usecases"

	uteam "sportlink/api/application/team/usecases"
	"sportlink/api/infrastructure/config"
	iaccount "sportlink/api/infrastructure/persistence/account"
	imatchoffer "sportlink/api/infrastructure/persistence/matchoffer"
	imatchrequest "sportlink/api/infrastructure/persistence/matchrequest"
	iplayer "sportlink/api/infrastructure/persistence/player"
	iteam "sportlink/api/infrastructure/persistence/team"

	"github.com/gin-gonic/gin"

	caccount "sportlink/api/infrastructure/rest/account"
	cauth "sportlink/api/infrastructure/rest/auth"
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
	accountRepository := iaccount.NewRepository(dynamoDbClient, "SportLinkCore")
	playerRepository := iplayer.NewDynamoDBRepository(dynamoDbClient, "SportLinkCore")
	teamRepository := iteam.NewRepository(dynamoDbClient, "SportLinkCore")
	matchOfferRepository := imatchoffer.NewRepository(dynamoDbClient, "SportLinkCore")
	matchRequestRepository := imatchrequest.NewRepository(dynamoDbClient, "SportLinkCore")
	matchRepository := imatch.NewRepository(dynamoDbClient, "SportLinkCore")

	// Account Use Cases
	findAccount := uaccount.NewFindAccountUC(accountRepository)

	// Player Use Cases
	createPlayer := uplayer.NewCreatePlayerUC(playerRepository)

	// Team Use Cases
	createTeam := uteam.NewCreateTeamUC(playerRepository, teamRepository)
	retrieveTeam := uteam.NewRetrieveTeamUC(teamRepository)
	findTeam := uteam.NewFindTeamUC(teamRepository)
	updateTeam := uteam.NewUpdateTeamUC(teamRepository)

	// Match Offer Use Cases
	createMatchOffer := umatchoffer.NewCreateMatchOfferUC(matchOfferRepository, teamRepository)
	findMatchOffers := umatchoffer.NewFindMatchOfferUC(matchOfferRepository)
	retrieveMatchOffer := umatchoffer.NewRetrieveMatchOfferUC(matchOfferRepository)
	deleteMatchOffer := umatchoffer.NewDeleteMatchOfferUC(matchOfferRepository)

	// Match Use Cases
	findMatches := umatch.NewFindMatchesUC(matchRepository)

	// Match Request Use Cases
	createMatchRequest := umatchrequest.NewCreateMatchRequestUC(matchRequestRepository, matchOfferRepository)
	findMatchRequests := umatchrequest.NewFindMatchRequestsUC(matchRequestRepository)
	updateMatchRequestStatus := umatchrequest.NewUpdateMatchRequestStatusUC(matchRequestRepository)
	acceptMatchRequest := umatchrequest.NewAcceptMatchRequestUC(matchRepository, matchRequestRepository, matchOfferRepository)

	// Auth Use Cases
	googleVerifier := authservice.NewGoogleTokenVerifier(cfg.AuthCfg.GoogleClientID)
	jwtService := authservice.NewJWTService(cfg.AuthCfg.JWTSecret)
	googleAuth := uauth.NewGoogleAuthUC(googleVerifier, accountRepository, jwtService)

	// Controllers
	accountController := caccount.NewController(findAccount)
	authController := cauth.NewController(googleAuth, customValidator)
	router.POST("/auth/google", authController.GoogleAuth)

	playerController := cplayer.NewController(&createPlayer, customValidator)
	router.POST("/player", playerController.CreatePlayer)

	teamController := cteam.NewController(createTeam, retrieveTeam, findTeam, findTeam, updateTeam, customValidator)
	router.GET("/account", accountController.Find)
	router.GET("/account/:account_id", accountController.Retrieve)

	router.POST("/account/:account_id/team", teamController.CreateTeam)
	router.GET("/account/:account_id/team", teamController.ListAccountTeams)
	router.GET("/sport/:sport/team/:team", teamController.RetrieveTeam)
	router.GET("/sport/:sport/team", teamController.FindTeam)
	router.PATCH("/sport/:sport/team/:team", teamController.UpdateTeam)

	matchOfferController := cmatchoffer.NewController(createMatchOffer, findMatchOffers, retrieveMatchOffer, deleteMatchOffer, customValidator)
	router.POST("/account/:account_id/match-offer", matchOfferController.CreateMatchOffer)
	router.GET("/match-offer", matchOfferController.FindMatchOffers)
	router.GET("/match-offer/:offer_id", matchOfferController.RetrieveMatchOffer)
	router.GET("/account/:account_id/match-offer", matchOfferController.FindAccountMatchOffers)
	router.DELETE("/account/:account_id/match-offer/:offer_id", matchOfferController.DeleteMatchOffer)

	matchRequestController := cmatchrequest.NewController(createMatchRequest, findMatchRequests, updateMatchRequestStatus, acceptMatchRequest, customValidator)
	router.POST("/account/:account_id/match-offer/:offer_id/match-request", matchRequestController.CreateMatchRequest)
	router.GET("/account/:account_id/match-request", matchRequestController.FindMatchRequests)
	router.PATCH("/account/:account_id/match-request/:request_id", matchRequestController.UpdateMatchRequestStatus)
	router.POST("/account/:account_id/match-request/:request_id/accept", matchRequestController.AcceptMatchRequest)

	matchController := cmatch.NewController(findMatches)
	router.GET("/account/:account_id/match", matchController.FindMatches)

	monitoring.RegisterMetricsRoute(router)
}
