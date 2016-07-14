// message
package control

import (
	"fmt"
)

type Message struct {
	receiver  string
	payload   string
	sender    string
	timestamp int
}

func sendMessage() {
	fmt.Println("Hello World!")
}

func receiveMessage() {
	fmt.Println("Hello World!")
}
