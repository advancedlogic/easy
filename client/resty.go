package client

import (
	"crypto/tls"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"net/http"
)

type Resty struct {
	Url         string
	QueryParams map[string]string
	Headers     map[string]string
	Cookies     map[string]string
	AuthToken   string
	Body        string
	Username    string
	Password    string
	pem         string
	key         string
}

func WithUrl(url string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if url != "" {
			r := client.(*Resty)
			r.Url = url
			return nil
		}
		return errors.New("url cannot be empty")
	}
}

func AddQueryParam(key, value string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if key != "" && value != "" {
			r := client.(*Resty)
			r.QueryParams[key] = value
			return nil
		}
		return errors.New("key and value cannot be empty")
	}
}

func AddHeader(key, value string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if key != "" && value != "" {
			r := client.(*Resty)
			r.Headers[key] = value
			return nil
		}
		return errors.New("key and value cannot be empty")
	}
}

func AddCookie(key, value string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if key != "" && value != "" {
			r := client.(*Resty)
			r.Cookies[key] = value
			return nil
		}
		return errors.New("key and value cannot be empty")
	}
}

func WithAuthToken(token string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if token != "" {
			r := client.(*Resty)
			r.AuthToken = token
			return nil
		}
		return errors.New("token cannot be empty")
	}
}

func WithBody(body string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if body != "" {
			r := client.(*Resty)
			r.Body = body
			return nil
		}
		return errors.New("body cannot be empty")
	}
}

func WithBasicAuthentication(username, password string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if username != "" && password != "" {
			r := client.(*Resty)
			r.Username = username
			r.Password = password
			return nil
		}
		return errors.New("username and password cannot be empty")
	}
}

func WithX509Certificate(pem, key string) interfaces.ClientOption {
	return func(client interfaces.Client) error {
		if pem != "" && key != "" {
			r := client.(*Resty)
			r.key = key
			r.pem = pem
			return nil
		}
		return errors.New("pem and key cannot be empty")
	}
}

func New(options ...interfaces.ClientOption) (*Resty, error) {
	r := &Resty{}
	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Resty) render() (*resty.Request, error) {
	client := resty.New()
	if len(r.Cookies) > 0 {
		for key, value := range r.Cookies {
			client.SetCookie(&http.Cookie{
				Name:  key,
				Value: value,
			})
		}
	}

	if r.pem != "" && r.key != "" {
		cert, err := tls.LoadX509KeyPair(r.pem, r.key)
		if err != nil {
			return nil, err
		}
		client.SetCertificates(cert)
	}

	request := client.R()
	if len(r.Headers) > 0 {
		request.SetHeaders(r.Headers)
	}
	if len(r.QueryParams) > 0 {
		request.SetQueryParams(r.QueryParams)
	}

	if r.AuthToken != "" {
		request.SetAuthToken(r.AuthToken)
	}
	if r.Body != "" {
		request.SetBody(r.Body)
	}
	if r.Username != "" && r.Password != "" {
		request.SetBasicAuth(r.Username, r.Password)
	}
	return request, nil
}

func (r *Resty) GET(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Get(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) POST(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Post(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) PUT(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Put(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) DELETE(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Delete(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) HEAD(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Head(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}

func (r *Resty) OPTIONS(h interface{}) error {
	request, err := r.render()
	if err != nil {
		return err
	}
	response, err := request.Options(r.Url)
	if err != nil {
		return err
	}
	handler := h.(func(response *resty.Response) error)
	return handler(response)
}
