package interfaces

type Transport interface {
	Run() error
	Stop() error

	Handler(string, string, interface{}) error
	Middleware(interface{}) error
	StaticFilesFolder(string, string) error
}

type TransportOption func(Transport) error
