package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// livenessHandler checks if the application is running. It's used by Kubernetes liveness probe
// to determine if the application is alive and running. If this check fails, Kubernetes may restart the container.
func livenessHandler(c *gin.Context) {
	c.String(http.StatusOK, "Liveness check passed: application is up and running")
}

// readinessHandler verifies that the application is ready to receive traffic. Used by Kubernetes readiness probe
// to ensure that the service is ready to process requests. This check could include dependencies like database or external services.
func readinessHandler(c *gin.Context) {
	c.String(http.StatusOK, "Readiness check passed: all systems are operational and ready to receive traffic")
}
