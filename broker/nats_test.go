package broker

import (
	"github.com/nats-io/go-nats"
	"gopkg.in/go-playground/assert.v1"
	"sync"
	"testing"
)

func TestNewNats(t *testing.T) {
	n, _ := NewNats()
	assert.NotEqual(t, n, nil)
}

func TestNats_Connect(t *testing.T) {
	n, _ := NewNats()
	err := n.Connect()
	defer n.Close()
	assert.Equal(t, err, nil)
}

func TestNats_Close(t *testing.T) {
	n, _ := NewNats()
	err := n.Connect()
	assert.Equal(t, err, nil)
	err = n.Close()
	assert.Equal(t, err, nil)
}

func TestNats_PublishSubscribe(t *testing.T) {
	wg := sync.WaitGroup{}
	n, _ := NewNats()
	err := n.Connect()
	assert.Equal(t, err, nil)
	defer n.Close()
	err = n.Subscribe("test", func(msg *nats.Msg) {
		assert.Equal(t, string(msg.Data), []byte("test"))
		wg.Done()
	})
	assert.Equal(t, err, nil)
	err = n.Publish("test", "test")
	assert.Equal(t, err, nil)
	wg.Wait()
}
