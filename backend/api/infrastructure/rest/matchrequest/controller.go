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
	AcceptMatchRequest(c *gin.Context)
}

type DefaultController struct {
	createMatchRequestUC       *usecases.CreateMatchRequestUC
	findMatchRequestsUC        *usecases.FindMatchRequestsUC
	updateMatchRequestStatusUC *usecases.UpdateMatchRequestStatusUC
	acceptMatchRequestUC       *usecases.AcceptMatchRequestUC
	validator                  *validator.Validate
}

func NewController(
	createMatchRequestUC *usecases.CreateMatchRequestUC,
	findMatchRequestsUC *usecases.FindMatchRequestsUC,
	updateMatchRequestStatusUC *usecases.UpdateMatchRequestStatusUC,
	acceptMatchRequestUC *usecases.AcceptMatchRequestUC,
	validator *validator.Validate,
) Controller {
	return &DefaultController{
		createMatchRequestUC:       createMatchRequestUC,
		findMatchRequestsUC:        findMatchRequestsUC,
		updateMatchRequestStatusUC: updateMatchRequestStatusUC,
		acceptMatchRequestUC:       acceptMatchRequestUC,
		validator:                  validator,
	}
}
