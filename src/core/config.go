// config
package core

import (
	"fmt"
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
	hosts []string
}

type httpConfig struct {
	Port string
}

var onceConfig sync.Once
var config Config

func NewConfig(file string) (*Config, error) {
	var initError error
	onceConfig.Do(func() {

		if _, err := toml.DecodeFile(file, &config); err != nil {
			fmt.Println("Error ", err)
			initError = err
		} else {
			fmt.Println("User is ", config)
		}
	})
	return &config, initError
}
