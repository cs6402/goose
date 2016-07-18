// message
package model

type Message struct {
	Receiver        string `json:"R"`
	Content         string `json:"C"`
	Sender          string `json:"S"`
	SenderMessageId string `json:"SMI"`
	Timestamp       int64  `json:"TS"`
	Type            int    `json:"T"`
}

type MessageWithId struct {
	Id      string `json:"I"`
	Payload string `json:"P"`
}
