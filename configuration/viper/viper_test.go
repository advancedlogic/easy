package viper

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNewViperConfiguration(t *testing.T) {
	c, _ := New(WithName("test"))
	err := c.Open("../../assets")
	assert.Equal(t, err, nil)
	assert.Equal(t, c.GetString("name"), "test")
	assert.Equal(t, c.GetStringOrDefault("test", "default"), "default")
}
