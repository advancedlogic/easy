package ledis

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestLedis_Init(t *testing.T) {
	ledis, err := New(AddEndpoints("localhost:6379"))
	assert.Equal(t, err, nil)
	assert.Equal(t, ledis.Init(), nil)
	assert.Equal(t, ledis.Close(), nil)
}

func TestLedis_Set(t *testing.T) {
	ledis, err := New(AddEndpoints("localhost:6379"))
	assert.Equal(t, err, nil)
	assert.Equal(t, ledis.Init(), nil)
	assert.Equal(t, ledis.Set("test", "test"), nil)
	assert.Equal(t, ledis.Close(), nil)
}

func TestLedis_Exists(t *testing.T) {
	ledis, err := New(AddEndpoints("localhost:6379"))
	assert.Equal(t, err, nil)
	assert.Equal(t, ledis.Init(), nil)
	exists, _ := ledis.Exists("test")
	assert.Equal(t, exists, true)
	assert.Equal(t, ledis.Close(), nil)
}

func TestLedis_Get(t *testing.T) {
	ledis, err := New(AddEndpoints("localhost:6379"))
	assert.Equal(t, err, nil)
	assert.Equal(t, ledis.Init(), nil)
	value, _ := ledis.Get("test")
	assert.Equal(t, value.(string), "test")
	assert.Equal(t, ledis.Close(), nil)
}

func TestLedis_Keys(t *testing.T) {
	ledis, err := New(AddEndpoints("localhost:6379"))
	assert.Equal(t, err, nil)
	assert.Equal(t, ledis.Init(), nil)
	keys, _ := ledis.Keys()
	assert.Equal(t, len(keys.([]string)) > 0, true)
	assert.Equal(t, ledis.Close(), nil)
}

func TestLedis_Delete(t *testing.T) {
	ledis, err := New(AddEndpoints("localhost:6379"))
	assert.Equal(t, err, nil)
	assert.Equal(t, ledis.Init(), nil)
	assert.Equal(t, ledis.Delete("test"), nil)
	assert.Equal(t, ledis.Close(), nil)
}
