package matchoffer

import (
	"sportlink/api/application"
	"sportlink/api/application/matchoffer/usecases"
	domainmatch "sportlink/api/domain/match"
	"sportlink/api/domain/matchoffer"
	"sportlink/api/infrastructure/rest/matchoffer/parser"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Controller interface {
	CreateMatchOffer(c *gin.Context)
	FindMatchOffers(c *gin.Context)
	FindAccountMatchOffers(c *gin.Context)
	RetrieveMatchOffer(c *gin.Context)
	DeleteMatchOffer(c *gin.Context)
	ConfirmMatchOffer(c *gin.Context)
}

type DefaultController struct {
	createMatchOfferUC   application.UseCase[matchoffer.Entity, matchoffer.Entity]
	findMatchOffersUC    application.UseCase[matchoffer.DomainQuery, usecases.FindMatchOfferResult]
	retrieveMatchOfferUC *usecases.RetrieveMatchOfferUC
	deleteMatchOfferUC   *usecases.DeleteMatchOfferUC
	confirmMatchOfferUC  application.UseCase[usecases.ConfirmMatchOfferInput, domainmatch.Entity]
	validator            *validator.Validate
	queryParser          parser.QueryParser
}

func NewController(
	createMatchOfferUC application.UseCase[matchoffer.Entity, matchoffer.Entity],
	findMatchOffersUC application.UseCase[matchoffer.DomainQuery, usecases.FindMatchOfferResult],
	retrieveMatchOfferUC *usecases.RetrieveMatchOfferUC,
	deleteMatchOfferUC *usecases.DeleteMatchOfferUC,
	confirmMatchOfferUC application.UseCase[usecases.ConfirmMatchOfferInput, domainmatch.Entity],
	validator *validator.Validate,
) Controller {
	return NewControllerWithParser(createMatchOfferUC, findMatchOffersUC, retrieveMatchOfferUC, deleteMatchOfferUC, confirmMatchOfferUC, validator, nil)
}

// NewControllerWithParser creates a controller with an optional query parser.
// If queryParser is nil, a new DefaultQueryParser will be created.
// This allows for dependency injection in tests.
func NewControllerWithParser(
	createMatchOfferUC application.UseCase[matchoffer.Entity, matchoffer.Entity],
	findMatchOffersUC application.UseCase[matchoffer.DomainQuery, usecases.FindMatchOfferResult],
	retrieveMatchOfferUC *usecases.RetrieveMatchOfferUC,
	deleteMatchOfferUC *usecases.DeleteMatchOfferUC,
	confirmMatchOfferUC application.UseCase[usecases.ConfirmMatchOfferInput, domainmatch.Entity],
	validator *validator.Validate,
	queryParser parser.QueryParser,
) Controller {
	if queryParser == nil {
		queryParser = parser.NewQueryParser()
	}
	return &DefaultController{
		createMatchOfferUC:   createMatchOfferUC,
		findMatchOffersUC:    findMatchOffersUC,
		retrieveMatchOfferUC: retrieveMatchOfferUC,
		deleteMatchOfferUC:   deleteMatchOfferUC,
		confirmMatchOfferUC:  confirmMatchOfferUC,
		validator:            validator,
		queryParser:          queryParser,
	}
}
