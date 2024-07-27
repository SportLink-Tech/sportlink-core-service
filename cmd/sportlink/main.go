package main

import (
	"github.com/gin-gonic/gin"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest"
)

func main() {
	r := gin.Default()

	// Register the error.go handling middleware
	r.Use(middleware.ErrorHandler())

	// Define routes
	rest.Routes(r)

	// Start the server
	r.Run()
}
