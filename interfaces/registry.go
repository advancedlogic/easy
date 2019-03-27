package interfaces

type Registry interface {
	Register() error
}

type RegistryOption func(Registry) error
