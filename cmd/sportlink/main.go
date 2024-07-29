package main

import (
	"github.com/gin-gonic/gin"
	"sportlink/api/infrastructure/middleware"
	"sportlink/api/infrastructure/rest"
)

func main() {
	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	rest.Routes(r)

	r.Run() // Inicia el servidor
}
