// cache
package main

import (
	"log"
	"sync"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type RedisCache struct {
	O        interface{}
	instance *pool.Pool
}

var once sync.Once

func NewCache() (*RedisCache, error) {

	var instance *RedisCache
	var initError error
	once.Do(func() {

		df := func(network, addr string) (*redis.Client, error) {
			client, err := redis.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			if err = client.Cmd("AUTH", "").Err; err != nil {
				client.Close()
				return nil, err
			}
			return client, nil
		}
		p, err := pool.NewCustom("tcp", "", 10, df)

		//		client, err := redis.Dial("tcp", "")
		if err != nil {
			log.Println("Cache init error", err)
			initError = err
		} else {
			instance = &RedisCache{instance: p}
		}
	})
	return instance, initError
}
