package matchrequest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	CreateMatchRequest(c *gin.Context)
	FindMatchRequests(c *gin.Context)
	FindSentMatchRequests(c *gin.Context)
	UpdateMatchRequestStatus(c *gin.Context)
}

type DefaultController struct {
	createMatchRequestUC       CreateMatchRequestUseCase
	findMatchRequestsUC        FindMatchRequestsUseCase
	findSentMatchRequestsUC    FindSentMatchRequestsUseCase
	updateMatchRequestStatusUC UpdateMatchRequestStatusUseCase
	validator                  *validator.Validate
}

func NewController(
	createMatchRequestUC CreateMatchRequestUseCase,
	findMatchRequestsUC FindMatchRequestsUseCase,
	findSentMatchRequestsUC FindSentMatchRequestsUseCase,
	updateMatchRequestStatusUC UpdateMatchRequestStatusUseCase,
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
