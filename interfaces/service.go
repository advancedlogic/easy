package interfaces

//Basic Microservice interface
// ID()  set/get unique Microservice identifier
// Name() set/get the Microservice name
// Run()
type Service interface {
	ID() string
	Name() string

	//Init(...ServiceOption)
	//Run(...interface{}) error
	//Stop(...interface{}) error
	//
	//Handle(string, string, interface{}) error
	//Subscribe(string, interface{}) error
}
