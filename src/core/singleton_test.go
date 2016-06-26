// singleton_test
package core

import (
	"fmt"
	"sync"
	"testing"
)

type singleton struct {
}

var instanceZ *singleton
var onceZ sync.Once

func GetInstance() *singleton {
	onceZ.Do(func() {
		instanceZ = &singleton{}
		fmt.Println("ffff", instanceZ)
	})
	return instanceZ
}

func estSingleton(t *testing.T) {
	z := GetInstance()
	fmt.Printf("%p", z)
	fmt.Println("")
	y := GetInstance()
	fmt.Printf("%p", y)
	fmt.Println("")
	fmt.Println("Hello World!", *z, *y)
}
