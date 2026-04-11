package match

import (
	"sportlink/api/application"
	"sportlink/api/application/match/usecases"
	domainmatch "sportlink/api/domain/match"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	FindMatches(c *gin.Context)
}

type DefaultController struct {
	findMatchesUC application.UseCase[usecases.FindMatchesInput, []domainmatch.Entity]
}

func NewController(findMatchesUC application.UseCase[usecases.FindMatchesInput, []domainmatch.Entity]) Controller {
	return &DefaultController{findMatchesUC: findMatchesUC}
}
