package main

import (
	"github.com/advancedlogic/easy/easy"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	microservice, err := easy.NewEasy(
		easy.WithDefaultRegistry(),
		easy.WithDefaultTransport())

	if err != nil {
		microservice.Fatal(err)
	}

	if err := microservice.Transport().Handler("get", "/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}); err != nil {
		microservice.Fatal(err)
	}

	microservice.Run()
	if err != nil {
		microservice.Fatal(err)
	}
}
