package consul

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/advancedlogic/easy/interfaces"
	"github.com/sirupsen/logrus"

	api "github.com/hashicorp/consul/api"
)

func WithID(id string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if id != "" {
			c := i.(*Consul)
			c.id = id
			return nil
		}
		return errors.New("ID cannot be empty")
	}
}

func WithName(name string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if name != "" {
			c := i.(*Consul)
			c.name = name
			return nil
		}
		return errors.New("name cannot be empty")
	}
}

func WithAddress(address string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if address != "" {
			c := i.(*Consul)
			c.address = address
			return nil
		}
		return errors.New("address cannot be empty")
	}
}

func WithPort(port int) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if port > 0 {
			c := i.(*Consul)
			c.port = port
			return nil
		}

		return errors.New("port cannot be zero")
	}
}

func WithUsername(username string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if username != "" {
			c := i.(*Consul)
			c.username = username
			return nil
		}
		return errors.New("username cannot be empty")
	}
}

func WithPassword(password string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if password != "" {
			c := i.(*Consul)
			c.password = password
			return nil
		}
		return errors.New("password cannot be empty")
	}
}

func WithCredentials(username, password string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if username != "" && password != "" {
			c := i.(*Consul)
			c.username = username
			c.password = password
			return nil
		}
		return errors.New("username/password cannot be empty")
	}
}

func WithInterval(interval string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if interval != "" {
			c := i.(*Consul)
			c.interval = interval
			return nil
		}
		return errors.New("interval cannot be empty")
	}
}

func WithTimeout(timeout string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if timeout != "" {
			c := i.(*Consul)
			c.timeout = timeout
			return nil
		}
		return errors.New("timeout cannot be empty")
	}
}

func WithHealthEndpoint(endpoint string) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		if endpoint != "" {
			c := i.(*Consul)
			c.healthEndpoint = endpoint
			return nil
		}
		return errors.New("endpoint cannot be empty")
	}
}

func WithLogger(logger *logrus.Logger) interfaces.RegistryOption {
	return func(i interfaces.Registry) error {
		c := i.(*Consul)
		return c.WithLogger(logger)
	}
}

type Consul struct {
	id             string
	name           string
	address        string
	port           int
	username       string
	password       string
	interval       string
	timeout        string
	healthEndpoint string
	*logrus.Logger
}

func New(options ...interfaces.RegistryOption) (*Consul, error) {
	//Default values first
	c := &Consul{
		id:             "default",
		name:           "default",
		address:        "localhost:8500",
		port:           8080,
		username:       "",
		password:       "",
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
	config.Address = c.address
	if c.username != "" && c.password != "" {
		config.HttpAuth = &api.HttpBasicAuth{
			Username: c.username,
			Password: c.password,
		}
	}
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

func (c *Consul) WithPort(port int) error {
	if port > -1 {
		c.port = port
		return nil
	}

	return errors.New("port must be positive")
}
