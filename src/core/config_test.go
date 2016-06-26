// config
package core

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {

	config, err := NewConfig("../config.toml")
	if err != nil {
		t.Errorf(err.Error())
	} else {
		z := config.RedisConfig
		ptr := &z
		fmt.Printf("%p", &ptr)
		fmt.Println("", reflect.TypeOf(z))
		vc := reflect.TypeOf(*config)
		//		fmt.Println(vc)

		for i := 0; i < vc.NumField(); i++ {
			fmt.Println("S", vc.Field(i))
		}

	}
}
func TestConfig2(t *testing.T) {

	config, err := NewConfig("../config.toml")

	if err != nil {
		t.Errorf(err.Error())
	} else {
		z := config.RedisConfig
		ptr := &z
		fmt.Printf("%p", &ptr)
		fmt.Println("")
		vc := reflect.TypeOf(*config)
		//		fmt.Println(vc)

		for i := 0; i < vc.NumField(); i++ {
			fmt.Println("S", vc.Field(i))
		}

	}
}
