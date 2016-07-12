// config
package core

import (
	"encoding/json"
	"flag"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RedisConfig    *redisConfig
	AwsConfig      *awsConfig
	DatabaseConfig *databaseConfig
	HttpConfig     *httpConfig
}

type redisConfig struct {
	Auth string
	Url  string
}
type awsConfig struct {
	AccessKey string
	SecretKey string
}

type databaseConfig struct {
	Hosts []string
}

type httpConfig struct {
	Port string
}

var onceConfig sync.Once
var config Config

func Get() *Config {
	return &config
}

func NewConfig() (*Config, error) {
	var initError error
	onceConfig.Do(func() {
		path := flag.String("config", "config.toml", "input config file path, or using default path: config.toml ")
		if _, err := toml.DecodeFile(*path, &config); err != nil {
			initError = err
		} else {
			body, _ := json.Marshal(&config)
			log.Println("Loading config succeed. Config:", *path, string(body))

		}
	})
	return &config, initError
}
