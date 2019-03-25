package interfaces

type Cache interface {
	Set(string, interface{}) error
	Get(string) (interface{}, error)
	Keys() (interface{}, error)
}
