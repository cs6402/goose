package repository

import (
	"core"
	"log"
	. "model"
	"time"
)

var begin time.Time

func init() {
	t, err := time.Parse("2006-01-02 15:04", "2006-01-02 15:04")
	if err != nil {
		panic(err)
	}
	begin = t
}

func AddMessage(msg *Message, payload string, ttl int) (result error) {
	session := core.NewCassandraConn()
	if err := session.Query(`INSERT INTO message (owner, id, payload) VALUES (?, now(), ?) USING TTL ?`,
		msg.Receiver, payload, ttl).Exec(); err != nil {
		result = err
	}
	return
}
func GetMessages(receiver string, last string, limit int) ([]*MessageWithId, error) {
	session := core.NewCassandraConn()
	iter := session.Query(`SELECT id, payload FROM message WHERE owner = ? and id >= ? LIMIT ?`,
		receiver, last, limit).Iter()
	result := make([]*MessageWithId, iter.NumRows())
	var id string
	var payload string
	index := 0
	for iter.Scan(&id, &payload) {
		result[index] = &MessageWithId{id, payload}
		index++
		log.Println(id, payload)
	}

	return result, iter.Close()
}

func GetMessagesFromBeginning(receiver string, limit int) ([]*MessageWithId, error) {
	session := core.NewCassandraConn()
	iter := session.Query(`SELECT id, payload FROM message WHERE owner = ? and id >= maxTimeuuid(?) LIMIT ?`,
		receiver, begin, limit).Iter()
	result := make([]*MessageWithId, iter.NumRows(), iter.NumRows())
	var id string
	var payload string
	log.Println(len(result))
	log.Println(iter.NumRows())
	index := 0
	for iter.Scan(&id, &payload) {
		result[index] = &MessageWithId{id, payload}
		index++
		log.Println(id, payload)
	}

	return result, iter.Close()
}
