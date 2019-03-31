package interfaces

type Cache interface {
	Init() error
	Close() error

	Set(string, interface{}) error
	Get(string) (interface{}, error)
	Exists(...string) (bool, error)
	Keys() (interface{}, error)
	Delete(...string) error
}

type CacheOption func(Cache) error
