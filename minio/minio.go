package minio

import (
	"errors"
	"github.com/minio/minio-go"
	"io/ioutil"
	"strings"
)

type Option func(*Minio) error

type Minio struct {
	location  string
	bucket    string
	endpoint  string
	accessKey string
	secretKey string
}

func WithLocation(location string) Option {
	return func(m *Minio) error {
		if location != "" {
			m.location = location
			return nil
		}
		return errors.New("location cannot be empty")
	}
}

func WithBucket(bucket string) Option {
	return func(m *Minio) error {
		if bucket != "" {
			m.bucket = bucket
			return nil
		}
		return errors.New("location cannot be empty")
	}
}

func WithEndpoint(endpoint string) Option {
	return func(m *Minio) error {
		if endpoint != "" {
			m.endpoint = endpoint
			return nil
		}
		return errors.New("endpoint cannot be empty")
	}
}

func WithAccessKey(accessKey string) Option {
	return func(m *Minio) error {
		if accessKey != "" {
			m.accessKey = accessKey
			return nil
		}
		return errors.New("access key cannot be empty")
	}
}

func WithSecretKey(secretKey string) Option {
	return func(m *Minio) error {
		if secretKey != "" {
			m.secretKey = secretKey
			return nil
		}
		return errors.New("secret key cannot be empty")
	}
}

func WithCredentials(accessKey, secretKey string) Option {
	return func(m *Minio) error {
		if accessKey != "" && secretKey != "" {
			m.accessKey = accessKey
			m.secretKey = secretKey
			return nil
		}
		return errors.New("access or secret key cannot be empty")
	}
}

func NewMinio(options ...Option) (*Minio, error) {
	m := &Minio{
		location: "default",
		bucket:   "default",
		endpoint: "localhost:9000",
	}

	for _, option := range options {
		if err := option(m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Minio) Create(key string, data interface{}) error {
	reader := strings.NewReader(data.(string))
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return err
	}
	_, err = client.PutObject(m.bucket, key, reader, -1, minio.PutObjectOptions{
		ContentType: "plain/txt",
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Minio) Read(key string) (interface{}, error) {
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return "", err
	}

	reader, err := client.GetObject(m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	if value, err := ioutil.ReadAll(reader); err == nil {
		return string(value), nil
	} else {
		return nil, err
	}
}

func (m *Minio) Update(key string, data interface{}) error {
	return m.Create(key, data)
}

func (m *Minio) Delete(key string) error {
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return err
	}

	return client.RemoveObject(m.bucket, key)
}

func (m *Minio) List(params ...interface{}) ([]interface{}, error) {
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return nil, err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	values := make([]interface{}, 0)
	for value := range client.ListObjectsV2(m.bucket, "", true, doneCh) {
		values = append(values, value)
	}
	return values, nil
}
