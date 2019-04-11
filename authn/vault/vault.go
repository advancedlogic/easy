package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advancedlogic/easy/commons"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/hashicorp/vault/api"
	"time"
)

type VaultUser struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Timestamp int64    `json:"timestamp"`
	Groups    []string `json:"groups"`
	Enabled   bool     `json:"enabled"`
}

type Vault struct {
	token               string
	servers             []string
	timeout             time.Duration
	skipTLSVerification bool
}

func New(options ...interfaces.AuthNOption) (*Vault, error) {
	v := &Vault{
		token:               "",
		servers:             make([]string, 0),
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

func WithToken(token string) interfaces.AuthNOption {
	return func(i interfaces.AuthN) error {
		if token != "" {
			vault := i.(*Vault)
			vault.token = token
			return nil
		}
		return errors.New("token cannot be nil")
	}
}

func WithServers(servers ...string) interfaces.AuthNOption {
	return func(i interfaces.AuthN) error {
		if len(servers) > 0 {
			vault := i.(*Vault)
			for _, server := range servers {
				vault.servers = append(vault.servers, server)
			}
			return nil
		}
		return errors.New("at least 1 server must be specified")
	}
}

func SkipTLS(skipTLSVerification bool) interfaces.AuthNOption {
	return func(i interfaces.AuthN) error {
		vault := i.(*Vault)
		vault.skipTLSVerification = skipTLSVerification
		return nil
	}
}

func (v *Vault) conenct() (*api.Client, error) {
	config := &api.Config{
		Address: v.servers[0],
	}
	if err := config.ConfigureTLS(&api.TLSConfig{
		Insecure: v.skipTLSVerification,
	}); err != nil {
		return nil, err
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(v.token)
	return client, nil
}

func (v *Vault) close(client *api.Client) error {
	return nil
}

func (v *Vault) Register(username, password string) (interface{}, error) {
	epassword, err := commons.HashAndSalt(password)
	if err != nil {
		return nil, err
	}
	user := VaultUser{
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

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	_, err = client.Logical().Write(fmt.Sprintf("/cubbyhole/%s", username), map[string]interface{}{
		"user": string(jsonUser),
	})
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

	secret, err := client.Logical().Read(fmt.Sprintf("/cubbyhole/%s", username))
	if err != nil {
		return nil, err
	}
	jsonUser := secret.Data["user"].(string)
	var user VaultUser
	err = json.Unmarshal([]byte(jsonUser), &user)
	if err != nil {
		return nil, err
	}
	if !commons.ComparePasswords(user.Password, []byte(password)) {
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

	_, err = client.Logical().Delete(fmt.Sprintf("/cubbyhole/%s", username))
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) Reset(username, password string) (interface{}, error) {
	return v.Register(username, password)
}
