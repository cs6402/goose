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

type ChatList struct {
	Owner string `cql:"owner"`
	// 1: chat, 2: secret
	ChatType int `cql:"type"`
	// key: chat_id, value: target
	ChatTarget map[string]string `cql:"chat"`
	// key: chat_id, value: last receive mid
	LastReceive map[string]string `cql:"last"`
}
