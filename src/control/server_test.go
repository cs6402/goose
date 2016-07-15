// server_test
package control

import (
	"fmt"
)

func TestRoute() {
	fmt.Println("Hello World!")
}

func call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f = reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
