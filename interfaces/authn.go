package interfaces

type AuthN interface {
	Register(string, string) (interface{}, error)
	Login(string, string) (interface{}, error)
	Logout(string) error
	Delete(string) error
	Reset(string) (interface{}, error)
}
