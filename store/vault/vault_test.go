package vault

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestVault_Create(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		WithNamespace("cubbyhole"),
		SkipTLSVerification(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	err = vault.Create("test", map[string]interface{}{
		"test": "test",
	})
	assert.Equal(t, err, nil)
}

func TestVault_Read(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		WithNamespace("cubbyhole"),
		SkipTLSVerification(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}

	err = vault.Create("test1", map[string]interface{}{
		"test1": "test1",
	})

	err = vault.Create("test2", map[string]interface{}{
		"test2": "test2",
	})

	data, err := vault.Read("test")
	assert.Equal(t, err, nil)
	m := data.(map[string]interface{})
	for k, v := range m {
		switch x := v.(type) {
		case string:
			assert.Equal(t, k, x)
		default:
			println(k, v)
		}
	}
}

func TestVault_List(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		WithNamespace("cubbyhole"),
		SkipTLSVerification(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	_, err = vault.List()
	assert.Equal(t, err, nil)
}

func TestVault_Delete(t *testing.T) {
	vault, err := New(
		WithToken("s.71PDtRZlmMEZfh7G94C9H7Wo"),
		WithNamespace("cubbyhole"),
		SkipTLSVerification(true),
		WithServers("https://127.0.0.1:8200"))
	if err != nil {
		panic(err)
	}
	err = vault.Delete("test")
	err = vault.Delete("test1")
	err = vault.Delete("test2")
	assert.Equal(t, err, nil)
}
