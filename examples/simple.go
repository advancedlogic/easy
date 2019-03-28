package main

import (
	"github.com/advancedlogic/easy/commons"
	. "github.com/advancedlogic/easy/easy"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	"net/http"
)

func main() {
	microservice, err := Default()

	if err != nil {
		panic(err)
	}

	if err := microservice.Broker().Subscribe("test", func(msg *nats.Msg) {
		microservice.Info(string(msg.Data))
	}); err != nil {
		microservice.Fatal(err)
	}

	endpointRequest := microservice.Configuration().GetStringOrDefault("endpoint.request", "ping")
	endpointResponse := microservice.Configuration().GetStringOrDefault("endpoint.response", "pong")

	if err := microservice.Transport().Handler(commons.ModeGet, endpointRequest, func(c *gin.Context) {
		c.String(http.StatusOK, endpointResponse)
	}); err != nil {
		microservice.Fatal(err)
	}

	microservice.Run()
}
