package vault

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestVault_Register(t *testing.T) {
	vault, err := New(
		WithToken("s.UA3hRV9lOjTKa8kojapKNW62"),
		SkipTLS(true),
		WithServers("http://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	_, err = vault.Register("test", "test")
	assert.Equal(t, err, nil)
}

func TestVault_Login(t *testing.T) {
	vault, err := New(
		WithToken("s.UA3hRV9lOjTKa8kojapKNW62"),
		SkipTLS(true),
		WithServers("http://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	_, err = vault.Login("test", "test")
	assert.Equal(t, err, nil)
}

func TestVault_Delete(t *testing.T) {
	vault, err := New(
		WithToken("s.UA3hRV9lOjTKa8kojapKNW62"),
		SkipTLS(true),
		WithServers("http://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	err = vault.Delete("test")
	assert.Equal(t, err, nil)
}
