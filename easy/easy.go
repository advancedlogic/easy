package easy

import (
	"errors"
	"github.com/advancedlogic/easy/authn/fs"
	"github.com/advancedlogic/easy/broker/nats"
	"github.com/advancedlogic/easy/commons"
	"github.com/advancedlogic/easy/configuration/viper"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/advancedlogic/easy/registry/consul"
	"github.com/advancedlogic/easy/transport/rest"
	"github.com/ankit-arora/go-utils/go-shutdown-hook"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"net/http"
	"plugin"
)

type Option func(*Easy) error

type Easy struct {
	id        string
	name      string
	isRunning bool

	registry      interfaces.Registry
	transport     interfaces.Transport
	broker        interfaces.Broker
	client        interfaces.Client
	store         interfaces.Store
	processor     interfaces.Processor
	configuration interfaces.Configuration
	authn         interfaces.AuthN
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
			return errors.New("name cannot be empty")
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
		r, err := consul.New(consul.WithLogger(easy.Logger))
		if err != nil {
			return err
		}
		easy.registry = r
		return nil
	}
}

func WithDefaultTransport() Option {
	return func(easy *Easy) error {
		t, err := rest.New(rest.WithLogger(easy.Logger))
		if err != nil {
			return err
		}
		easy.transport = t
		return nil
	}
}

func WithDefaultBroker() Option {
	return func(easy *Easy) error {
		b, err := nats.New(nats.WithLogger(easy.Logger))
		if err != nil {
			return err
		}
		easy.broker = b
		return nil
	}
}

func WithDefaultConfiguration() Option {
	return func(easy *Easy) error {
		c, err := viper.New(
			viper.WithName(easy.name))
		if err != nil {
			return err
		}
		err = c.Open()
		if err != nil {
			return err
		}
		easy.configuration = c
		return nil

	}
}

func WithDefaultAuthN(folder string) Option {
	return func(easy *Easy) error {
		if folder != "" {
			c, err := fs.New(
				fs.WithFolder(folder))
			if err != nil {
				return err
			}
			easy.authn = c
			return nil
		}
		return errors.New("folder cannot be empty")
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

func WithClient(client interfaces.Client) Option {
	return func(easy *Easy) error {
		if client != nil {
			easy.client = client
			return nil
		}
		return errors.New("client cannot be nil")
	}
}

func WithStore(store interfaces.Store) Option {
	return func(easy *Easy) error {
		if store != nil {
			easy.store = store
			return nil
		}
		return errors.New("store cannot be nil")
	}
}

func WithProcessor(processor interfaces.Processor) Option {
	return func(easy *Easy) error {
		if processor != nil {
			err := processor.Init(easy)
			if err != nil {
				return err
			}
			easy.processor = processor
			return nil
		}
		return errors.New("processor cannot be nil")
	}
}

func WithPlugin(lib, name string) Option {
	return func(easy *Easy) error {
		if lib != "" && name != "" {
			plug, err := plugin.Open(lib)
			if err != nil {
				return err
			}
			symPlugin, err := plug.Lookup(name)
			if err != nil {
				return err
			}
			var processor interfaces.Processor
			processor, ok := symPlugin.(interfaces.Processor)
			if !ok {
				return errors.New("unexpected type from module symbol")
			}
			easy.processor = processor
			return nil
		}
		return errors.New("lib and name cannot be empty")
	}
}

func WithConfiguration(configuration interfaces.Configuration) Option {
	return func(easy *Easy) error {
		if configuration != nil {
			easy.configuration = configuration
			return nil
		}
		return errors.New("configuration cannot be nil")
	}
}

func WithLocalConfiguration() Option {
	return func(easy *Easy) error {
		if easy.name != "" {
			conf, err := viper.New(
				viper.WithName(easy.name),
				viper.WithLogger(easy.Logger))
			if err != nil {
				return err
			}
			if err := conf.Open(); err != nil {
				return err
			}
			easy.configuration = conf
		}
		return errors.New("name cannot be empty")
	}
}

func WithRemoteConfiguration(provider, uri string) Option {
	return func(easy *Easy) error {
		if provider != "" && uri != "" {
			conf, err := viper.New(
				viper.WithName(easy.name),
				viper.WithProvider(provider),
				viper.WithURI(uri),
				viper.WithLogger(easy.Logger))
			if err != nil {
				return nil
			}
			easy.configuration = conf
		}

		return errors.New("provider and uri cannot be empty")
	}
}

//NewEasy create a new µs according to the passed options
//WithID: default random
//WithName: default "default"
func New(options ...Option) (*Easy, error) {
	easy := &Easy{
		id:     uuid.New().String(),
		Logger: logrus.New(),
		name:   "default",
	}

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	easy.Formatter = formatter

	for _, option := range options {
		err := option(easy)
		if err != nil {
			return nil, err
		}
	}

	logLevel := "info"
	if easy.configuration != nil {
		if err := easy.configuration.Open(); err != nil {
			return nil, err
		}
		logLevel = easy.configuration.GetStringOrDefault("log.level", "info")
		if timestamp := easy.configuration.GetStringOrDefault("log.timestamp", ""); timestamp != "" {
			formatter.TimestampFormat = timestamp
		}
	}
	switch logLevel {
	case "debug":
		easy.Level = logrus.DebugLevel
	case "info":
		easy.Level = logrus.InfoLevel
	case "warn":
		easy.Level = logrus.WarnLevel
	case "error":
		easy.Level = logrus.ErrorLevel
	default:
		easy.Level = logrus.InfoLevel
	}

	return easy, nil
}

func Default(options ...Option) (*Easy, error) {
	microservice, err := New(
		WithDefaultConfiguration(),
		WithDefaultRegistry(),
		WithDefaultBroker(),
		WithDefaultTransport())
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		if err := option(microservice); err != nil {
			return nil, err
		}
	}
	return microservice, nil
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

func (easy *Easy) Client() interfaces.Client {
	return easy.client
}

func (easy *Easy) Store() interfaces.Store {
	return easy.store
}

func (easy *Easy) Configuration() interfaces.Configuration {
	return easy.configuration
}

func (easy *Easy) Processor() interfaces.Processor {
	return easy.processor
}

func (easy *Easy) AuthN() interfaces.AuthN {
	return easy.authn
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

	if easy.authn != nil {
		easy.Info("authn setup")

		register := func(c *gin.Context) {
			var user fs.User
			err := c.BindJSON(&user)
			if err != nil {
				c.String(http.StatusBadGateway, err.Error())
				return
			}
			response, err := easy.authn.Register(user.Username, user.Password)
			if err != nil {
				c.String(http.StatusBadGateway, err.Error())
				return
			}
			c.JSON(http.StatusOK, response)
		}

		login := func(c *gin.Context) {
			var user fs.User
			err := c.BindJSON(&user)
			if err != nil {
				c.String(http.StatusBadGateway, err.Error())
				return
			}
			response, err := easy.authn.Login(user.Username, user.Password)
			if err != nil {
				c.String(http.StatusBadGateway, err.Error())
			}
			c.JSON(http.StatusOK, response)
		}

		logout := func(c *gin.Context) {
			username := c.Param(":username")
			err := easy.authn.Logout(username)
			if err != nil {
				c.String(http.StatusBadGateway, err.Error())
				return
			}
			c.String(http.StatusOK, "")
		}

		if err := easy.transport.Handler("post", "/register", register); err != nil {
			easy.Fatal(err)
		}

		if err := easy.transport.Handler("post", "/login", login); err != nil {
			easy.Fatal(err)
		}

		if err := easy.transport.Handler("post", "/logout", logout); err != nil {
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

func (easy *Easy) Handler(mode, route string, handler interface{}) error {
	return easy.transport.Handler(mode, route, handler)
}

func (easy *Easy) GET(route string, handler interface{}) error {
	return easy.transport.Handler(commons.ModeGet, route, handler)
}

func (easy *Easy) POST(route string, handler interface{}) error {
	return easy.transport.Handler(commons.ModePost, route, handler)
}

func (easy *Easy) PUT(route string, handler interface{}) error {
	return easy.transport.Handler(commons.ModePut, route, handler)
}

func (easy *Easy) DELETE(route string, handler interface{}) error {
	return easy.transport.Handler(commons.ModeDelete, route, handler)
}

func (easy *Easy) Subscribe(endpoint string, handler interface{}) error {
	return easy.broker.Subscribe(endpoint, handler)
}

func (easy *Easy) Unsubscribe(endpoint string) error {
	return easy.broker.Unsubscribe(endpoint)
}

func (easy *Easy) Publish(endpoint string, msg interface{}) error {
	return easy.broker.Publish(endpoint, msg)
}

func (easy *Easy) Info(message interface{}) {
	easy.Logger.Info(message)
}

func (easy *Easy) Warn(message interface{}) {
	easy.Logger.Warn(message)
}

func (easy *Easy) Error(message interface{}) {
	easy.Logger.Error(message)
}

func (easy *Easy) Fatal(message interface{}) {
	easy.Logger.Fatal(message)
}

func (easy *Easy) Debug(message interface{}) {
	easy.Logger.Debug(message)
}
