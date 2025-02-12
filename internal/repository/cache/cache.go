package cache

import "time"


type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, bool)
	Delete(key string) error
	Clear() error
}
