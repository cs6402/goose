package repository

import (
	"core"
	"log"
	. "model"
)

func AddMessage(msg *Message, payload string, ttl int) (result error) {
	session := core.NewCassandraWConn()
	if err := session.Query(`INSERT INTO message (owner, id, payload) VALUES (?, now(), ?) USING TTL ?`,
		msg.Receiver, payload, ttl).Exec(); err != nil {
		result = err
	}
	return
}
func GetMessages(receiver string, last string, limit int) ([]*MessageWithId, error) {
	session := core.NewCassandraRConn()
	iter := session.Query(`SELECT id, payload FROM message WHERE owner = ? and id > ? ORDER BY id ASC LIMIT ?`,
		receiver, last, limit).Iter()
	result := make([]*MessageWithId, iter.NumRows())
	var id string
	var payload string
	index := iter.NumRows() - 1
	for iter.Scan(&id, &payload) {
		result[index] = &MessageWithId{id, payload}
		index--
		//		log.Println(id, payload)
	}

	return result, iter.Close()
}

func GetMessagesFromBeginning(receiver string, limit int) ([]*MessageWithId, error) {
	session := core.NewCassandraRConn()
	iter := session.Query(`SELECT id, payload FROM message WHERE owner = ? ORDER BY id ASC LIMIT ?`,
		receiver, limit).Iter()
	result := make([]*MessageWithId, iter.NumRows(), iter.NumRows())
	var id string
	var payload string
	log.Println(len(result))
	log.Println(iter.NumRows())
	index := iter.NumRows() - 1
	for iter.Scan(&id, &payload) {
		result[index] = &MessageWithId{id, payload}
		index--
		//		log.Println(id, payload)
	}

	return result, iter.Close()
}

func ConfirmMessages(receiver string, start string, end string, limit int, counts int) ([]*MessageWithId, error) {
	session := core.NewCassandraCConn()
	var dbCounts int
	var endId string
	err := session.Query(`SELECT COUNT(*), id FROM message WHERE owner = ? and id >= ? and id <= ?`,
		receiver, start, end).Scan(&dbCounts, &endId)

	if err != nil {
		return nil, err
	}
	if dbCounts == counts && endId == end {
		// normal
		//		wsession := core.NewCassandraWConn()
		//		wsession.Query(``)
		log.Println("Confirm messages owner:", receiver, " from:", start, " to:", end, "counts:", counts)
		return nil, nil
	} else {
		// unexpected
		iter := session.Query(`SELECT id, payload FROM message WHERE owner = ? and id >= ? and id <= ?`,
			receiver, start, end).Iter()
		result := make([]*MessageWithId, iter.NumRows())
		var id string
		var payload string
		index := 0
		for iter.Scan(&id, &payload) {
			result[index] = &MessageWithId{id, payload}
			index++
			//			log.Println("Confirm ", id, payload)
		}
		if length := len(result); length != 0 && length == counts && result[0].Id == end && result[length-1].Id == start {
			return nil, iter.Close()
		}

		return result, iter.Close()
	}

}
