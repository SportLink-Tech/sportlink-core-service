package auth

import (
	"net/http"
	"sportlink/api/application/auth/request"
	"sportlink/api/application/errors"

	"github.com/gin-gonic/gin"
)

// GoogleAuth handles POST /auth/google
// Verifies the Google id_token, creates the account if needed, and returns a session cookie.
func (sc *DefaultController) GoogleAuth(c *gin.Context) {
	var req request.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.InvalidRequestFormat())
		return
	}

	if err := sc.validator.Struct(req); err != nil {
		c.Error(errors.RequestValidationFailed(err.Error()))
		return
	}

	result, err := sc.googleAuthUC.Invoke(c.Request.Context(), req.IDToken)
	if err != nil {
		c.Error(errors.Unauthorized(err.Error()))
		return
	}

	// httpOnly=true: JS cannot read the cookie (XSS safe)
	// secure=false: allow HTTP in local dev; set to true in production
	c.SetCookie("token", result.JWTToken, 7*24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"account_id": result.AccountID})
}
