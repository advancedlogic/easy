package interfaces

type Cache interface {
	Init() error
	Close() error

	Put(string, interface{}) error
	Take(string) (interface{}, error)
	Exists(...string) (bool, error)
	Keys() (interface{}, error)
	Delete(...string) error

	Set(string) error
	IsMember(string) (bool, error)
}

type CacheOption func(Cache) error
