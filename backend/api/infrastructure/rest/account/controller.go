package account

import (
	"github.com/gin-gonic/gin"

	"sportlink/api/application"
	"sportlink/api/application/account/usecases"
	"sportlink/api/domain/account"
)

type Controller interface {
	Retrieve(c *gin.Context)
	Find(c *gin.Context)
}

type DefaultController struct {
	findAccountUC application.UseCase[usecases.FindAccountInput, []account.Entity]
}

func NewController(findAccountUC application.UseCase[usecases.FindAccountInput, []account.Entity]) Controller {
	return &DefaultController{
		findAccountUC: findAccountUC,
	}
}
