package match

import (
	"sportlink/api/application"
	"sportlink/api/application/match/usecases"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	FindMatches(c *gin.Context)
}

type DefaultController struct {
	findMatchesUC application.UseCase[usecases.FindMatchesInput, []usecases.MatchWithOffer]
}

func NewController(findMatchesUC application.UseCase[usecases.FindMatchesInput, []usecases.MatchWithOffer]) Controller {
	return &DefaultController{findMatchesUC: findMatchesUC}
}
