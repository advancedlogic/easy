package registry

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"

	api "github.com/hashicorp/consul/api"
)

type Option func(*Consul) error

func WithID(id string) Option {
	return func(c *Consul) error {
		if id != "" {
			c.id = id
			return nil
		}
		return errors.New("ID cannot be empty")
	}
}

func WithName(name string) Option {
	return func(c *Consul) error {
		if name != "" {
			c.name = name
			return nil
		}
		return errors.New("name cannot be empty")
	}
}

func WithPort(port int) Option {
	return func(c *Consul) error {
		if port > 0 {
			c.port = port
			return nil
		}

		return errors.New("port cannot be zero")
	}
}

func WithInterval(interval string) Option {
	return func(c *Consul) error {
		if interval != "" {
			c.interval = interval
			return nil
		}
		return errors.New("interval cannot be empty")
	}
}

func WithTimeout(timeout string) Option {
	return func(c *Consul) error {
		if timeout != "" {
			c.timeout = timeout
			return nil
		}
		return errors.New("timeout cannot be empty")
	}
}

func WithHealthEndpoint(endpoint string) Option {
	return func(c *Consul) error {
		if endpoint != "" {
			c.healthEndpoint = endpoint
			return nil
		}
		return errors.New("endpoint cannot be empty")
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(c *Consul) error {
		return c.WithLogger(logger)
	}
}

type Consul struct {
	id             string
	name           string
	port           int
	interval       string
	timeout        string
	healthEndpoint string
	*logrus.Logger
}

func NewConsul(options ...Option) (*Consul, error) {
	//Default values first
	c := &Consul{
		id:             "default",
		name:           "default",
		port:           8080,
		interval:       "3s",
		timeout:        "5s",
		healthEndpoint: "",
		Logger:         logrus.New(),
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Consul) Register() error {
	hostname := func() string {
		hn, err := os.Hostname()
		if err != nil {
			log.Fatalln(err)
		}
		return hn
	}

	config := api.DefaultConfig()
	consul, err := api.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	registration := new(api.AgentServiceRegistration)
	registration.ID = c.id
	registration.Name = c.name
	address := hostname()
	registration.Address = address
	registration.Port = c.port
	if c.healthEndpoint != "" {
		registration.Check = new(api.AgentServiceCheck)
		registration.Check.HTTP = fmt.Sprintf("http://%s:%v/%s",
			address, c.port, c.healthEndpoint)
		registration.Check.Interval = c.interval
		registration.Check.Timeout = c.timeout
	}
	return consul.Agent().ServiceRegister(registration)
}

func (c *Consul) WithLogger(logger *logrus.Logger) error {
	if logger != nil {
		c.Logger = logger
		return nil
	}
	return errors.New("logger cannot be nil")
}
