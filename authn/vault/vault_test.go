package vault

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestVault_Register(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		SkipTLS(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	_, err = vault.Register("test", "test")
	assert.Equal(t, err, nil)
}

func TestVault_Login(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		SkipTLS(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	_, err = vault.Login("test", "test")
	assert.Equal(t, err, nil)
}

func TestVault_Delete(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		SkipTLS(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	err = vault.Delete("test")
	assert.Equal(t, err, nil)
}
