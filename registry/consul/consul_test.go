package consul

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConsul(t *testing.T) {
	c, _ := New()
	assert.NotEqual(t, c, nil)
}

func TestWithID(t *testing.T) {
	c, _ := New(WithID("123"))
	assert.Equal(t, c.id, "123")
}

func TestWithName(t *testing.T) {
	c, _ := New(WithName("test"))
	assert.Equal(t, c.name, "test")
}

func TestWithPort(t *testing.T) {
	c, _ := New(WithPort(9090))
	assert.Equal(t, c.port, 9090)
}

func TestWithInterval(t *testing.T) {
	c, _ := New(WithInterval("1s"))
	assert.Equal(t, c.interval, "1s")
}

func TestWithTimeout(t *testing.T) {
	c, _ := New(WithTimeout("1s"))
	assert.Equal(t, c.timeout, "1s")
}

func TestWithHealthEndpoint(t *testing.T) {
	c, _ := New(WithHealthEndpoint("health"))
	assert.Equal(t, c.healthEndpoint, "health")
}

func TestConsul_Register(t *testing.T) {
	c, _ := New()
	err := c.Register()
	assert.Equal(t, err, nil)
}
