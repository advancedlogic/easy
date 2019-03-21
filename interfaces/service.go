package interfaces

import "github.com/advancedlogic/goms/pkg/interfaces"

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
	//
	//Handle(string, string, interface{}) error
	//Subscribe(string, interface{}) error

	Registry() interfaces.Registry
	Transport() interfaces.Transport
}
