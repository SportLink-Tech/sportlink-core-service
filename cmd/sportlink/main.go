package main

import (
	"github.com/gin-gonic/gin"
	"sportlink/api/infrastructure/rest"
)

func main() {
	r := gin.Default()
	rest.Routes(r)
	r.Run()
}
