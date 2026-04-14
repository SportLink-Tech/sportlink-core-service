package matchrequest

import (
	"sportlink/api/application"
	"sportlink/api/application/matchrequest/usecases"
	"sportlink/api/domain/matchrequest"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	CreateMatchRequest(c *gin.Context)
	FindMatchRequests(c *gin.Context)
	UpdateMatchRequestStatus(c *gin.Context)
	AcceptMatchRequest(c *gin.Context)
	CancelMatchRequest(c *gin.Context)
}

type DefaultController struct {
	createMatchRequestUC       *usecases.CreateMatchRequestUC
	findMatchRequestsUC        *usecases.FindMatchRequestsUC
	updateMatchRequestStatusUC *usecases.UpdateMatchRequestStatusUC
	acceptMatchRequestUC       application.UseCase[usecases.AcceptMatchRequestInput, matchrequest.Entity]
	cancelMatchRequestUC       application.UseCase[usecases.CancelMatchRequestInput, matchrequest.Entity]
	validator                  *validator.Validate
}

func NewController(
	createMatchRequestUC *usecases.CreateMatchRequestUC,
	findMatchRequestsUC *usecases.FindMatchRequestsUC,
	updateMatchRequestStatusUC *usecases.UpdateMatchRequestStatusUC,
	acceptMatchRequestUC application.UseCase[usecases.AcceptMatchRequestInput, matchrequest.Entity],
	cancelMatchRequestUC application.UseCase[usecases.CancelMatchRequestInput, matchrequest.Entity],
	validator *validator.Validate,
) Controller {
	return &DefaultController{
		createMatchRequestUC:       createMatchRequestUC,
		findMatchRequestsUC:        findMatchRequestsUC,
		updateMatchRequestStatusUC: updateMatchRequestStatusUC,
		acceptMatchRequestUC:       acceptMatchRequestUC,
		cancelMatchRequestUC:       cancelMatchRequestUC,
		validator:                  validator,
	}
}
