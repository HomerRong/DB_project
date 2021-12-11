package gredis

import (
	"github.com/go-redis/redis"
	settings "main/pkg/setting"
	"time"
)

var rc *redis.Client

// Setup Initialize the Redis instance
func Setup() {

	rc = redis.NewClient(&redis.Options{
		Addr:     settings.RedisSetting.Addr,
		Password: settings.RedisSetting.Password,
		DB:       settings.RedisSetting.DB,
		// Dual
		IdleTimeout: settings.RedisSetting.IdleTimeout,
	})

	if err := rc.Ping().Err(); err != nil {
		panic(err)
	}
}

// Set setting a key/value
func Set(key string, value interface{}, exp time.Duration) error {
	return rc.Set(key, value, exp).Err()
}

func Exist(key string) (bool, error) {
	val, err := rc.Exists(key).Result()
	return val > 0, err
}

func Get(key string) (string, error) {
	val, err := rc.Get(key).Result()
	return val, err
}

func Delete(key string) error {
	return rc.Del(key).Err()
}
