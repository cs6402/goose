// cache
package core

import (
	"log"
	"sync"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

var onceCache sync.Once
var cache *pool.Pool

func NewCache() *pool.Pool {
	onceCache.Do(func() {

		df := func(network, addr string) (*redis.Client, error) {
			client, err := redis.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			if err = client.Cmd("AUTH", config.RedisConfig.Auth).Err; err != nil {
				client.Close()
				return nil, err
			}
			return client, nil
		}
		p, err := pool.NewCustom("tcp", config.RedisConfig.Url, 10, df)

		if err != nil {
			log.Panicln("Cache init failed.", err)
		} else {
			cache = p
		}
	})
	return cache
}
