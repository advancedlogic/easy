package authn

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/easy/commons"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

type FS struct {
	folder string
}

type FSOption func(*FS) error

type FSUser struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Timestamp int64    `json:"timestamp"`
	Groups    []string `json:"groups"`
	Enabled   bool     `json:"enabled"`
}

func NewFSUser(username, password string) (*FSUser, error) {
	if username != "" && password != "" {
		epassword, err := commons.HashAndSalt(password)
		if err != nil {
			return nil, err
		}
		return &FSUser{
			Username:  username,
			Password:  epassword,
			Timestamp: time.Now().UnixNano(),
			Groups:    []string{"user"},
			Enabled:   true,
		}, nil
	}
	return nil, errors.New("username and password cannot be empty")
}

func WithFolder(folder string) FSOption {
	return func(fs *FS) error {
		if folder != "" {
			fs.folder = folder
			return nil
		}
		return errors.New("folder cannot be empty")
	}
}

func NewFS(options ...FSOption) (*FS, error) {
	fs := &FS{
		folder: "fs",
	}
	for _, option := range options {
		if err := option(fs); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func (f *FS) Register(username, password string) (interface{}, error) {
	if username != "" && password != "" {
		user, err := NewFSUser(username, password)
		if err != nil {
			return nil, err
		}
		jsonUser, err := json.Marshal(user)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", f.folder, username), jsonUser, 0644)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, errors.New("username and password cannot be empty")
}

func (f *FS) Login(username, password string) (interface{}, error) {
	if username != "" && password != "" {
		jsonUser, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", f.folder, username))
		if err != nil {
			return nil, err
		}
		var user FSUser
		err = json.Unmarshal(jsonUser, &user)
		if err != nil {
			return nil, err
		}

		epassword, err := commons.HashAndSalt(password)
		if err != nil {
			return nil, err
		}

		if epassword != user.Password {
			return nil, errors.New("wrong username or password")
		}
		return user, nil
	}
	return nil, errors.New("username and password cannot be empty")
}

func (f *FS) Logout(username string) error {
	if username != "" {
		return nil
	}
	return errors.New("username cannot be empty")
}

func (f *FS) Delete(username string) error {
	if username != "" {
		return nil
	}
	return errors.New("username cannot be empty")
}

func (f *FS) Reset(username, password string) (interface{}, error) {
	return f.Register(username, password)
}
