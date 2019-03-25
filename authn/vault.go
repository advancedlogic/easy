package authn

import (
	"encoding/json"
	"errors"
	"github.com/shoenig/vaultapi"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Timestamp int64    `json:"timestamp"`
	Groups    []string `json:"groups"`
	Enabled   bool     `json:"enabled"`
}

type Option func(*Vault) error

type Vault struct {
	token               string
	servers             []string
	timeout             time.Duration
	skipTLSVerification bool
}

func NewVault(options ...Option) (*Vault, error) {
	v := &Vault{
		token:               "",
		servers:             []string{"http://localhost:8200"},
		timeout:             10 * time.Second,
		skipTLSVerification: true,
	}

	for _, option := range options {
		err := option(v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func WithToken(token string) Option {
	return func(vault *Vault) error {
		if token != "" {
			vault.token = token
			return nil
		}
		return errors.New("token cannot be nil")
	}
}

func WithServers(servers ...string) Option {
	return func(vault *Vault) error {
		if len(servers) > 0 {
			for _, server := range servers {
				vault.servers = append(vault.servers, server)
			}
			return nil
		}
		return errors.New("at least 1 server must be specified")
	}
}

func (v *Vault) conenct() (vaultapi.Client, error) {
	options := vaultapi.ClientOptions{
		Servers:             v.servers,
		HTTPTimeout:         v.timeout,
		SkipTLSVerification: v.skipTLSVerification,
	}

	tokener := vaultapi.NewStaticToken(v.token)
	return vaultapi.New(options, tokener)
}

func (v *Vault) close(client vaultapi.Client) error {
	return client.StepDown()
}

func (v *Vault) Register(username, password string) (interface{}, error) {
	epassword, err := v.hashAndSalt(password)
	if err != nil {
		return nil, err
	}
	user := User{
		Username:  username,
		Password:  epassword,
		Groups:    []string{"user"},
		Timestamp: time.Now().UnixNano(),
		Enabled:   true,
	}

	client, err := v.conenct()
	if err != nil {
		return nil, err
	}
	defer v.close(client)
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = client.Put(username, string(jsonUser))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (v *Vault) Login(username, password string) (interface{}, error) {
	client, err := v.conenct()
	if err != nil {
		return nil, err
	}
	defer v.close(client)
	jsonUser, err := client.Get(username)
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal([]byte(jsonUser), &user)
	if err != nil {
		return nil, err
	}
	epassword, err := v.hashAndSalt(password)
	if err != nil {
		return nil, err
	}
	if user.Password != epassword {
		return nil, errors.New("wrong username or password")
	}
	return user, nil
}

func (v *Vault) Logout(username string) error {
	return nil
}

func (v *Vault) Delete(username string) error {
	client, err := v.conenct()
	if err != nil {
		return err
	}
	defer v.close(client)
	err = client.Delete(username)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) Reset(username, password string) (interface{}, error) {
	return v.Register(username, password)
}

func (v *Vault) hashAndSalt(password string) (string, error) {
	bpassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bpassword, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
