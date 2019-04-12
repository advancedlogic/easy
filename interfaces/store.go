package interfaces

type Store interface {
	Create(string, interface{}) error
	Read(string) (interface{}, error)
	Update(string, interface{}) error
	Delete(string) error
	List(...interface{}) (interface{}, error)
}

type StoreOption func(Store) error
