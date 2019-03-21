package main

import (
	. "github.com/advancedlogic/easy/easy"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	"net/http"
)

func main() {
	microservice, err := NewEasy(
		WithDefaultRegistry(),
		WithDefaultBroker(),
		WithDefaultTransport())

	if err != nil {
		microservice.Fatal(err)
	}

	if err := microservice.Subscribe("test", func(msg *nats.Msg) {

	}); err != nil {
		microservice.Fatal(err)
	}

	if err := microservice.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}); err != nil {
		microservice.Fatal(err)
	}

	microservice.Run()
	if err != nil {
		microservice.Fatal(err)
	}
}
