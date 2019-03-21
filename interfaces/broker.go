package interfaces

import "github.com/nats-io/go-nats"

type Broker interface {
	Run() error
	Endpoint() string
	Connect() error
	Publish(string, interface{}) error
	Subscribe(string, nats.MsgHandler) error
	Unsubscribe(string) error
	Close() error
}
