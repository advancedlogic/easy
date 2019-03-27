package interfaces

type AuthN interface {
	Register(string, string) (interface{}, error)
	Login(string, string) (interface{}, error)
	Logout(string) error
	Delete(string) error
	Reset(string, string) (interface{}, error)
}

type AuthNOption func(AuthN) error
