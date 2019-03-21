package easy

import (
	"errors"
	"github.com/advancedlogic/easy/broker"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/advancedlogic/easy/registry"
	"github.com/advancedlogic/easy/transport"
	"github.com/ankit-arora/go-utils/go-shutdown-hook"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Option func(*Easy) error

type Easy struct {
	id        string
	name      string
	isRunning bool

	registry  interfaces.Registry
	transport interfaces.Transport
	broker    interfaces.Broker

	*logrus.Logger
}

//WithID(id string) set the id of the µs
func WithID(id string) Option {
	return func(easy *Easy) error {
		if id == "" {
			return errors.New("ID cannot be empty")
		}
		easy.id = id
		return nil
	}
}

//WithName(name string) set the id of the µs
func WithName(name string) Option {
	return func(easy *Easy) error {
		if name == "" {
			return errors.New("Name cannot be empty")
		}
		easy.name = name
		return nil
	}
}

func WithRegistry(registry interfaces.Registry) Option {
	return func(easy *Easy) error {
		if registry != nil {
			easy.registry = registry
			return nil
		}
		return errors.New("registry cannot be nil")
	}
}

func WithDefaultRegistry() Option {
	return func(easy *Easy) error {
		r, err := registry.NewConsul(registry.WithLogger(easy.Logger))
		if err != nil {
			return err
		}
		easy.registry = r
		return nil
	}
}

func WithDefaultTransport() Option {
	return func(easy *Easy) error {
		t, err := transport.NewRest(transport.WithLogger(easy.Logger))
		if err != nil {
			return err
		}
		easy.transport = t
		return nil
	}
}

func WithDefaultBroker() Option {
	return func(easy *Easy) error {
		b, err := broker.NewNats(broker.WithLogger(easy.Logger))
		if err != nil {
			return err
		}
		easy.broker = b
		return nil
	}
}

func WithHandler(mode, route string, handler interface{}) Option {
	return func(easy *Easy) error {
		return easy.transport.Handler(mode, route, handler)
	}
}

func WithMiddleware(middleware interface{}) Option {
	return func(easy *Easy) error {
		return easy.transport.Middleware(middleware)
	}
}

func WithStaticFilesFolder(route, folder string) Option {
	return func(easy *Easy) error {
		return easy.transport.StaticFilesFolder(route, folder)
	}
}

func WithTransport(transport interfaces.Transport) Option {
	return func(easy *Easy) error {
		if transport != nil {
			easy.transport = transport
			return nil
		}
		return errors.New("transport cannot be nil")
	}
}

func WithBroker(broker interfaces.Broker) Option {
	return func(easy *Easy) error {
		if broker != nil {
			easy.broker = broker
			return nil
		}
		return errors.New("broker cannot be nil")
	}
}

//NewEasy create a new µs according to the passed options
//WithID: default random
//WithName: default "default"
func NewEasy(options ...Option) (*Easy, error) {
	//Default values
	easy := &Easy{
		id:     uuid.New().String(),
		name:   "default",
		Logger: logrus.New(),
	}
	for _, option := range options {
		err := option(easy)
		if err != nil {
			return nil, err
		}
	}
	return easy, nil
}

//ID() return the µs' ID
//Part of Service interface implementation
func (easy *Easy) ID() string {
	return easy.id
}

//Name() return the µs' name
//Part of Service interface implementation
func (easy *Easy) Name() string {
	return easy.name
}

//Registry() return the µs' registry
func (easy *Easy) Registry() interfaces.Registry {
	return easy.registry
}

//Transport() return the µs' transport
func (easy *Easy) Transport() interfaces.Transport {
	return easy.transport
}

func (easy *Easy) Broker() interfaces.Broker {
	return easy.broker
}

func (easy *Easy) Run() {
	go_shutdown_hook.ADD(func() {
		easy.Stop()
		easy.Warn("Goodbye and thanks for all the fish")
	})
	if easy.registry != nil {
		easy.Info("registry setup")
		err := easy.registry.Register()
		if err != nil {
			easy.Fatal(err)
		}
	}
	if easy.broker != nil {
		easy.Info("broker setup")
		err := easy.broker.Run()
		if err != nil {
			easy.Fatal(err)
		}
	}
	if easy.transport != nil {
		easy.Info("transport setup")
		err := easy.transport.Run()
		if err != nil {
			easy.Fatal(err)
		}
	}
	easy.isRunning = true

	go_shutdown_hook.Wait()
}

func (easy *Easy) Stop() {
	if easy.broker != nil {
		if err := easy.broker.Close(); err != nil {
			easy.Fatal(err)
		}
	}
	if easy.transport != nil {
		if err := easy.transport.Stop(); err != nil {
			easy.Fatal(err)
		}
	}
}

func (easy *Easy) IsRunning() bool {
	return easy.isRunning
}

func (easy *Easy) HookShutDown(fn func()) {
	go_shutdown_hook.ADD(fn)
}
