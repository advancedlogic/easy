package easy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEasy(t *testing.T) {
	easy, _ := NewEasy()
	assert.NotEqual(t, easy, nil)
}

func TestWithID(t *testing.T) {
	easy, _ := NewEasy(WithID("123"))
	assert.Equal(t, easy.id, "123")
}

func TestWithName(t *testing.T) {
	easy, _ := NewEasy(WithName("test"))
	assert.Equal(t, easy.name, "test")
}

func TestErrorInID(t *testing.T) {
	_, err := NewEasy(WithID(""))
	assert.NotEqual(t, err, nil)
}

func TestErrorInName(t *testing.T) {
	_, err := NewEasy(WithName(""))
	assert.NotEqual(t, err, nil)
}

func TestWithDefaultRegistry(t *testing.T) {
	easy, _ := NewEasy(WithDefaultRegistry())
	assert.NotEqual(t, easy, nil)
}
