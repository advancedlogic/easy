package ledis

import (
	"github.com/advancedlogic/easy/interfaces"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Ledis struct {
	collection    string
	endpoints     []string
	username      string
	password      string
	db            int
	clusterClient *redis.ClusterClient
	client        *redis.Client
}

func WithCollection(collection string) interfaces.CacheOption {
	return func(c interfaces.Cache) error {
		if collection != "" {
			ledis := c.(*Ledis)
			ledis.collection = collection
			return nil
		}
		return errors.New("collection cannot be empty")
	}
}

func WithPassowrd(password string) interfaces.CacheOption {
	return func(c interfaces.Cache) error {
		if password != "" {
			ledis := c.(*Ledis)
			ledis.password = password
			return nil
		}
		return errors.New("username and password cannot be empty")
	}
}

func WithDB(db int) interfaces.CacheOption {
	return func(c interfaces.Cache) error {
		if db > -1 {
			ledis := c.(*Ledis)
			ledis.db = db
			return nil
		}
		return errors.New("db must be >= 0")
	}
}

func AddEndpoints(endpoints ...string) interfaces.CacheOption {
	return func(c interfaces.Cache) error {
		for _, endpoint := range endpoints {
			if endpoint != "" {
				ledis := c.(*Ledis)
				ledis.endpoints = append(ledis.endpoints, endpoint)
				return nil
			}
		}
		return errors.New("endpoint cannot be empty")
	}
}

func New(options ...interfaces.CacheOption) (*Ledis, error) {
	ledis := &Ledis{
		endpoints: make([]string, 0),
	}
	for _, option := range options {
		if err := option(ledis); err != nil {
			return nil, err
		}
	}
	return ledis, nil
}

func (l *Ledis) Init() error {
	if len(l.endpoints) > 1 {
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    l.endpoints,
			Password: l.password,
		})
		status := clusterClient.Ping()
		if status.Err() != nil {
			return status.Err()
		}
		l.clusterClient = clusterClient
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     l.endpoints[0],
			Password: l.password,
			DB:       l.db,
		})
		l.client = client
	}
	return nil
}

func (l *Ledis) Close() error {
	if l.clusterClient != nil {
		return l.clusterClient.Close()
	}
	return nil
}

func (l *Ledis) Put(key string, value interface{}) error {
	var status *redis.StatusCmd
	if l.client != nil {
		status = l.client.Set(key, value, -1)
	} else {
		status = l.clusterClient.Set(key, value, -1)
	}
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (l *Ledis) Set(value string) error {
	var status *redis.IntCmd
	if l.client != nil {
		status = l.client.SAdd(l.collection, value)
	} else {
		status = l.clusterClient.SAdd(l.collection, value)
	}
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (l *Ledis) IsMember(value string) (bool, error) {
	var status *redis.BoolCmd
	if l.client != nil {
		status = l.client.SIsMember(l.collection, value)
	} else {
		status = l.clusterClient.SIsMember(l.collection, value)
	}
	if status.Err() != nil {
		return false, status.Err()
	}
	return status.Val(), nil
}

func (l *Ledis) Take(key string) (interface{}, error) {
	var status *redis.StringCmd
	if l.client != nil {
		status = l.client.Get(key)
	} else {
		status = l.clusterClient.Get(key)
	}
	if status.Err() != nil {
		return nil, status.Err()
	}
	result := status.Val()
	return result, nil
}

func (l *Ledis) Exists(keys ...string) (bool, error) {
	var status *redis.IntCmd
	if l.client != nil {
		status = l.client.Exists(keys...)
	} else {
		status = l.clusterClient.Exists(keys...)
	}
	if status.Err() != nil {
		return false, status.Err()
	}
	return status.Val() > 0, nil
}

func (l *Ledis) Delete(keys ...string) error {
	var status *redis.IntCmd
	if l.client != nil {
		status = l.client.Del(keys...)
	} else {
		status = l.clusterClient.Del(keys...)
	}
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (l *Ledis) Keys() (interface{}, error) {
	var status *redis.StringSliceCmd
	if l.client != nil {
		status = l.client.Keys("*")
	} else {
		status = l.clusterClient.Keys("*")
	}
	if status.Err() != nil {
		return nil, status.Err()
	}
	return status.Val(), nil
}
