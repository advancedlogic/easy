package viper

import (
	"errors"
	"fmt"
	"time"

	"github.com/advancedlogic/easy/interfaces"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func WithName(name string) interfaces.ConfigurationOption {
	return func(i interfaces.Configuration) error {
		if name != "" {
			v := i.(*Viper)
			v.name = name
			return nil
		}

		return errors.New("provider cannot be empty")
	}
}

func WithProvider(provider string) interfaces.ConfigurationOption {
	return func(i interfaces.Configuration) error {
		if provider != "" {
			v := i.(*Viper)
			v.provider = provider
			return nil
		}

		return errors.New("provider cannot be empty")
	}
}

func WithURI(uri string) interfaces.ConfigurationOption {
	return func(i interfaces.Configuration) error {
		if uri != "" {
			v := i.(*Viper)
			v.uri = uri
			return nil
		}

		return errors.New("uri cannot be empty")
	}
}

func WithLogger(logger *logrus.Logger) interfaces.ConfigurationOption {
	return func(i interfaces.Configuration) error {
		if logger != nil {
			v := i.(*Viper)
			v.Logger = logger
			return nil
		}

		return errors.New("logger cannot be nil")
	}
}

type Viper struct {
	*viper.Viper
	*logrus.Logger

	name     string
	provider string
	uri      string
}

func New(options ...interfaces.ConfigurationOption) (*Viper, error) {
	v := &Viper{
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

func (v *Viper) Open(paths ...string) error {
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

func (v *Viper) Save(data interface{}) error {
	return nil
}

func (v *Viper) Set(key string, value interface{}) error {
	return nil
}

func (v *Viper) Get(key string) (interface{}, error) {
	return nil, nil
}

func (v *Viper) GetIntOrDefault(path string, defaultValue int) int {
	if value := v.GetInt(path); value != 0 {
		return value
	}
	return defaultValue
}

func (v *Viper) GetBoolOrDefault(path string, defaultValue bool) bool {
	if value := v.GetBool(path); value {
		return value
	}
	return defaultValue
}

func (v *Viper) GetFloat64OrDefault(path string, defaultValue float64) float64 {
	if value := v.GetFloat64(path); value != 0.0 {
		return value
	}
	return defaultValue
}

func (v *Viper) GetDurationOrDefault(path string, defaultValue time.Duration) time.Duration {
	if value := v.GetDuration(path); value != 0.0 {
		return value
	}
	return defaultValue
}

func (v *Viper) GetStringOrDefault(path string, defaultValue string) string {
	if value := v.GetString(path); value != "" {
		return value
	}
	return defaultValue
}

func (v *Viper) GetArrayOfStringsOrDefault(path string, defaultValue []string) []string {
	if value := v.GetStringSlice(path); value != nil {
		return value
	}
	return defaultValue
}
