package matchrequest

import (
	"sportlink/api/application/matchrequest/usecases"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	CreateMatchRequest(c *gin.Context)
	FindMatchRequests(c *gin.Context)
	UpdateMatchRequestStatus(c *gin.Context)
}

type DefaultController struct {
	createMatchRequestUC       *usecases.CreateMatchRequestUC
	findMatchRequestsUC        *usecases.FindMatchRequestsUC
	findSentMatchRequestsUC    *usecases.FindSentMatchRequestsUC
	updateMatchRequestStatusUC *usecases.UpdateMatchRequestStatusUC
	validator                  *validator.Validate
}

func NewController(
	createMatchRequestUC *usecases.CreateMatchRequestUC,
	findMatchRequestsUC *usecases.FindMatchRequestsUC,
	findSentMatchRequestsUC *usecases.FindSentMatchRequestsUC,
	updateMatchRequestStatusUC *usecases.UpdateMatchRequestStatusUC,
	validator *validator.Validate,
) Controller {
	return &DefaultController{
		createMatchRequestUC:       createMatchRequestUC,
		findMatchRequestsUC:        findMatchRequestsUC,
		findSentMatchRequestsUC:    findSentMatchRequestsUC,
		updateMatchRequestStatusUC: updateMatchRequestStatusUC,
		validator:                  validator,
	}
}
