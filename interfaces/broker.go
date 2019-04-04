package interfaces

type Broker interface {
	Run() error
	Endpoint() string
	Connect() error
	Publish(string, interface{}) error
	Subscribe(string, interface{}) error
	Unsubscribe(string) error
	Close() error
}

type BrokerOption func(Broker) error
