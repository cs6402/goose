// message
package control

import (
	"fmt"
)

type Message struct {
	Receiver        string `json:"R"`
	Content         string `json:"C"`
	Sender          string `json:"S"`
	SenderMessageId string `json:"SI"`
	Timestamp       int64  `json:"T"`
}

func sendMessage() {
	fmt.Println("Hello World!")
}

func receiveMessage() {
	fmt.Println("Hello World!")
}
