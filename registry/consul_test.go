package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConsul(t *testing.T) {
	c, _ := NewConsul()
	assert.NotEqual(t, c, nil)
}

func TestWithID(t *testing.T) {
	c, _ := NewConsul(WithID("123"))
	assert.Equal(t, c.id, "123")
}

func TestWithName(t *testing.T) {
	c, _ := NewConsul(WithName("test"))
	assert.Equal(t, c.name, "test")
}

func TestWithPort(t *testing.T) {
	c, _ := NewConsul(WithPort(9090))
	assert.Equal(t, c.port, 9090)
}

func TestWithInterval(t *testing.T) {
	c, _ := NewConsul(WithInterval("1s"))
	assert.Equal(t, c.interval, "1s")
}

func TestWithTimeout(t *testing.T) {
	c, _ := NewConsul(WithTimeout("1s"))
	assert.Equal(t, c.timeout, "1s")
}

func TestWithHealthEndpoint(t *testing.T) {
	c, _ := NewConsul(WithHealthEndpoint("health"))
	assert.Equal(t, c.healthEndpoint, "health")
}

func TestConsul_Register(t *testing.T) {
	c, _ := NewConsul()
	err := c.Register()
	assert.Equal(t, err, nil)
}
