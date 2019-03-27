package fs

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/easy/commons"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/pkg/errors"
	"io/ioutil"
	"time"
)

type FS struct {
	folder string
}

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Timestamp int64    `json:"timestamp"`
	Groups    []string `json:"groups"`
	Enabled   bool     `json:"enabled"`
}

func NewUser(username, password string) (*User, error) {
	if username != "" && password != "" {
		epassword, err := commons.HashAndSalt(password)
		if err != nil {
			return nil, err
		}
		return &User{
			Username:  username,
			Password:  epassword,
			Timestamp: time.Now().UnixNano(),
			Groups:    []string{"user"},
			Enabled:   true,
		}, nil
	}
	return nil, errors.New("username and password cannot be empty")
}

func WithFolder(folder string) interfaces.AuthNOption {
	return func(a interfaces.AuthN) error {
		if folder != "" {
			fs := a.(*FS)
			fs.folder = folder
			return nil
		}
		return errors.New("folder cannot be empty")
	}
}

func NewFS(options ...interfaces.AuthNOption) (*FS, error) {
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
		user, err := NewUser(username, password)
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
		user.Password = ""
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
		var user User
		err = json.Unmarshal(jsonUser, &user)
		if err != nil {
			return nil, err
		}

		if !commons.ComparePasswords(user.Password, []byte(password)) {
			return nil, errors.New("wrong username or password")
		}
		user.Password = ""
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
