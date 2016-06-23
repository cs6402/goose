// config
package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pelletier/go-toml"
)

//type Config struct {
//	redisConfig RedisConfig
//}

//type RedisConfig struct {
//	auth string
//	url  string
//}

func TestL(t *testing.T) {
	filename := "config.toml"
	//	var config Config
	config, err := toml.LoadFile(filename)
	if err != nil {
		fmt.Println("Error ", err.Error())
		t.Errorf("Reverse(%q) == %q, want %q")

	} else {
		// retrieve data directly
		user := config.Get("redis.url")
		fmt.Println("User is ", user, reflect.TypeOf(user))
		//		password := config.Get("postgres.password").(string)

		//		// or using an intermediate object
		//		configTree := config.Get("postgres").(*toml.TomlTree)
		//		user = configTree.Get("user").(string)
		//		password = configTree.Get("password").(string)
		//		fmt.Println("User is ", user, ". Password is ", password)

		//		// show where elements are in the file
		//		fmt.Println("User position: %v", configTree.GetPosition("user"))
		//		fmt.Println("Password position: %v", configTree.GetPosition("password"))

		//		// use a query to gather elements without walking the tree
		//		results, _ := config.Query("$..[user,password]")
		//		for ii, item := range results.Values() {
		//			fmt.Println("Query result %d: %v", ii, item)
		//		}
	}
}
