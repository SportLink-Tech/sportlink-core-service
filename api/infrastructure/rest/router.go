package rest

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	router.GET("/livez", livenessHandler)
	router.GET("/readyz", readinessHandler)

	RegisterMetricsRoute(router)
}
