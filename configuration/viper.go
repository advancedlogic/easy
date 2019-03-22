package configuration

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Option func(*ViperConfiguration) error

func WithName(name string) Option {
	return func(v *ViperConfiguration) error {
		if name != "" {
			v.name = name
			return nil
		}

		return errors.New("provider cannot be empty")
	}
}

func WithProvider(provider string) Option {
	return func(v *ViperConfiguration) error {
		if provider != "" {
			v.provider = provider
			return nil
		}

		return errors.New("provider cannot be empty")
	}
}

func WithURI(uri string) Option {
	return func(v *ViperConfiguration) error {
		if uri != "" {
			v.uri = uri
			return nil
		}

		return errors.New("uri cannot be empty")
	}
}

func WithConfigFile(configFile string) Option {
	return func(v *ViperConfiguration) error {
		if configFile != "" {
			v.local = configFile
			return nil
		}

		return errors.New("config file cannot be empty")
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(v *ViperConfiguration) error {
		if logger != nil {
			v.Logger = logger
			return nil
		}

		return errors.New("logger cannot be nil")
	}
}

type ViperConfiguration struct {
	*viper.Viper
	*logrus.Logger

	name     string
	provider string
	uri      string
	local    string
}

func NewViperConfiguration(options ...Option) (*ViperConfiguration, error) {
	v := &ViperConfiguration{
		Viper:  viper.New(),
		Logger: logrus.New(),
	}

	for _, option := range options {
		if err := option(v); err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (v *ViperConfiguration) Open(paths ...string) error {
	v.SetConfigName(v.name)

	if v.provider != "" && v.uri != "" {
		if err := v.AddRemoteProvider(v.provider, v.uri, v.name); err != nil {
			return err
		}
		if err := v.ReadRemoteConfig(); err != nil {
			return err
		}
	} else {
		v.AddConfigPath(fmt.Sprintf("/etc/%s/", v.name))
		v.AddConfigPath(fmt.Sprintf("$HOME/.%s", v.name))
		for _, path := range paths {
			v.AddConfigPath(path)
		}
		v.AddConfigPath(".")
		v.AutomaticEnv()
		if err := v.ReadInConfig(); err != nil {
			return err
		}
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		err := v.ReadInConfig()
		if err != nil {
			return
		}
	})

	return nil
}

func (v *ViperConfiguration) Save(data interface{}) error {
	return nil
}

func (v *ViperConfiguration) Set(key string, value interface{}) error {
	return nil
}

func (v *ViperConfiguration) Get(key string) (interface{}, error) {
	return nil, nil
}

func (v *ViperConfiguration) GetIntOrDefault(path string, defaultValue int) int {
	if value := v.GetInt(path); value != 0 {
		return value
	}
	return defaultValue
}

func (v *ViperConfiguration) GetBoolOrDefault(path string, defaultValue bool) bool {
	if value := v.GetBool(path); value {
		return value
	}
	return defaultValue
}

func (v *ViperConfiguration) GetFloat64OrDefault(path string, defaultValue float64) float64 {
	if value := v.GetFloat64(path); value != 0.0 {
		return value
	}
	return defaultValue
}

func (v *ViperConfiguration) GetDurationOrDefault(path string, defaultValue time.Duration) time.Duration {
	if value := v.GetDuration(path); value != 0.0 {
		return value
	}
	return defaultValue
}

func (v *ViperConfiguration) GetStringOrDefault(path string, defaultValue string) string {
	if value := v.GetString(path); value != "" {
		return value
	}
	return defaultValue
}

func (v *ViperConfiguration) Log() *logrus.Logger {
	return v.Logger
}
