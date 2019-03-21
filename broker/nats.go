package broker

import (
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

type Option func(*Nats) error

type Nats struct {
	endpoint            string
	conn                *nats.Conn
	userCredentialsPath string
	userJWT             string
	userNK              string
	*logrus.Logger
	handlers      map[string]func(*nats.Msg)
	subscriptions map[string]*nats.Subscription
}

func WithEndpoint(endpoint string) Option {
	return func(i *Nats) error {
		if endpoint != "" {
			i.endpoint = endpoint
			return nil
		}
		return errors.New("endpoint cannot be empty")
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(i *Nats) error {
		return i.WithLogger(logger)
	}
}

func NewNats(options ...Option) (*Nats, error) {
	n := &Nats{
		endpoint:      "localhost:4222",
		handlers:      make(map[string]func(*nats.Msg)),
		subscriptions: make(map[string]*nats.Subscription),
		Logger:        logrus.New(),
	}

	for _, option := range options {
		if err := option(n); err != nil {
			return nil, err
		}
	}

	return n, nil
}

func (n *Nats) Connect() error {
	var err error
	if conn, err := nats.Connect(n.endpoint); err == nil {
		n.conn = conn
		return nil
	}
	return err

}

func (n *Nats) Publish(topic string, message interface{}) error {
	return n.conn.Publish(topic, message.([]byte))
}

func (n *Nats) Subscribe(topic string, handler interface{}) error {
	n.handlers[topic] = handler.(func(*nats.Msg))
	return nil
}

func (n *Nats) Unsubscribe(topic string) error {
	if subscription, exists := n.subscriptions[topic]; exists {
		return subscription.Unsubscribe()
	}
	return errors.New(fmt.Sprintf("topic %s does not exist", topic))
}

func (n *Nats) Run() error {
	if err := n.Connect(); err != nil {
		return err
	}
	for topic, handler := range n.handlers {
		subscription, err := n.conn.Subscribe(topic, handler)
		if err != nil {
			return err
		}
		n.subscriptions[topic] = subscription
	}
	return nil
}

func (n *Nats) Endpoint() string {
	return n.endpoint
}

func (n *Nats) Close() error {
	if n.conn != nil {
		n.conn.Close()
		return nil
	}
	return errors.New("broker cannot be closed")
}

func (n *Nats) WithLogger(logger *logrus.Logger) error {
	if logger != nil {
		n.Logger = logger
		return nil
	}
	return errors.New("logger cannot be nil")
}