package easy

import (
	"errors"
	"github.com/google/uuid"
)

type Option func(*Easy) error

type Easy struct {
	id   string
	name string
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

//NewEasy create a new µs according to the passed options
//WithID: default random
//WithName: default "default"
func NewEasy(options ...Option) (*Easy, error) {
	//Default values
	easy := &Easy{
		id:   uuid.New().String(),
		name: "default",
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
