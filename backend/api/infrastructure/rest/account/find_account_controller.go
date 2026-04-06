package account

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sportlink/api/application/account/usecases"
	"sportlink/api/application/errors"
)

func (sc *DefaultController) Find(c *gin.Context) {
	accountID := strings.TrimSpace(c.Query("account_id"))
	email := strings.TrimSpace(c.Query("email"))

	if accountID == "" && email == "" {
		c.Error(errors.RequestValidationFailed("account_id or email query parameter is required"))
		return
	}
	if accountID != "" && email != "" {
		c.Error(errors.RequestValidationFailed("provide only one of account_id or email"))
		return
	}

	accs, err := sc.findAccountUC.Invoke(c.Request.Context(), usecases.FindAccountInput{
		AccountID: accountID,
		Email:     email,
	})
	if err != nil {
		c.Error(errors.UseCaseExecutionFailed(err.Error()))
		return
	}

	if accs == nil || len(*accs) == 0 {
		c.Error(errors.NotFound("No account found"))
		return
	}

	c.JSON(http.StatusOK, *accs)
}
