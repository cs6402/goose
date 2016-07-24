// message
package service

import (
	"log"
	. "model"
	"repository"
)

func SendMessage(msg *Message, payload string) {
	// panic recovery

	err := repository.AddMessage(msg, payload, 50)
	if err != nil {
		log.Println("send message failed!", msg.Sender, " to ", msg.Receiver, " body", payload)
	}
}

func ReceiveMessage(receiver string, last string) []*MessageWithId {
	// panic recovery
	result, err := repository.GetMessages(receiver, last, 50)
	if err != nil {
		log.Println("receive message failed", receiver, last)
	}
	return result
}
