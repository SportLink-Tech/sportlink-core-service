package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportlink/api/application/account/usecases"
	"sportlink/api/application/errors"
	"strings"
)

func (sc *DefaultController) Retrieve(c *gin.Context) {
	accountID := strings.TrimSpace(c.Param("account_id"))

	if accountID == "" {
		c.Error(errors.RequestValidationFailed("account_id is required"))
		return
	}

	accounts, err := sc.findAccountUC.Invoke(c.Request.Context(), usecases.FindAccountInput{
		AccountID: accountID,
	})

	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	if accounts == nil || len(*accounts) == 0 {
		c.Error(errors.NotFound("No account found"))
		return
	}
	
	c.JSON(http.StatusOK, (*accounts)[0])
}
