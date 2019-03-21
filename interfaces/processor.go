package interfaces

type Processor interface {
	Init(service Service) error
	Close() error
	Process(interface{}) (interface{}, error)
}
