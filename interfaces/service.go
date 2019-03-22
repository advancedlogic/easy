package interfaces

//Basic Microservice interface
// ID()  set/get unique Microservice identifier
// Name() set/get the Microservice name
// Run()
type Service interface {
	ID() string
	Name() string

	//Init(...ServiceOption)
	Run()
	Stop()
	IsRunning() bool
	HookShutDown(func())

	Registry() Registry
	Transport() Transport
	Broker() Broker
	Client() Client
	Store() Store

	//Transport Handler (rest) Helpers
	Handler(string, string, interface{}) error
	GET(string, interface{}) error
	POST(string, interface{}) error
	PUT(string, interface{}) error
	DELETE(string, interface{}) error

	//Broker Handler Helpers
	Subscribe(string, interface{}) error
	Unsubscribe(string) error
	Publish(string, interface{}) error
	//Store

	//Log
	Info(interface{})
	Warn(interface{})
	Error(interface{})
	Fatal(interface{})
	Debug(interface{})
}
