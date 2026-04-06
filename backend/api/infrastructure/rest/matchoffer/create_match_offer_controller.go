package matchoffer

import (
	"net/http"
	"sportlink/api/application/errors"
	"sportlink/api/application/matchoffer/mapper"
	"sportlink/api/application/matchoffer/request"
	restmapper "sportlink/api/infrastructure/rest/matchoffer/mapper"

	"github.com/gin-gonic/gin"
)

// CreateMatchOffer handles the POST request to create a match offer.
func (sc *DefaultController) CreateMatchOffer(c *gin.Context) {
	var newMatchOfferRequest request.NewMatchOfferRequest
	if err := c.ShouldBindJSON(&newMatchOfferRequest); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(newMatchOfferRequest); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	accountID := c.Param("account_id")
	input, err := mapper.CreationRequestToEntity(newMatchOfferRequest, accountID)
	if err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.createMatchOfferUC.Invoke(c.Request.Context(), input)
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	// Convert domain entity to response DTO
	responseDTO := restmapper.EntityToResponse(*result)

	c.JSON(http.StatusCreated, responseDTO)
}
