package broker

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestNewNats(t *testing.T) {
	n, _ := NewNats()
	assert.NotEqual(t, n, nil)
}
