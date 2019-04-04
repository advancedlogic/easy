package interfaces

type Client interface {
	GET(interface{}) error
	POST(interface{}) error
	PUT(interface{}) error
	DELETE(interface{}) error
	HEAD(interface{}) error
	OPTIONS(interface{}) error
}

type ClientOption func(Client) error
