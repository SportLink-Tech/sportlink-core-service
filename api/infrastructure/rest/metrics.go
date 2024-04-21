package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RegisterMetricsRoute registers the /metrics endpoint.
func RegisterMetricsRoute(router *gin.Engine) {
	h := promhttp.Handler()
	router.GET("/metrics", gin.WrapH(h))
}
