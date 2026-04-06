package auth

import (
	"sportlink/api/application/auth/usecases"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	GoogleAuth(c *gin.Context)
}

type DefaultController struct {
	googleAuthUC *usecases.GoogleAuthUC
	validator    *validator.Validate
}

func NewController(googleAuthUC *usecases.GoogleAuthUC, validator *validator.Validate) Controller {
	return &DefaultController{
		googleAuthUC: googleAuthUC,
		validator:    validator,
	}
}
