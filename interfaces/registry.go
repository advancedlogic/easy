package interfaces

type Registry interface {
	Register() error
	WithPort(port int) error
}

type RegistryOption func(Registry) error
