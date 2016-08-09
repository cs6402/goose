package repository

import (
	"core"
	"log"
	. "model"
)

const (
	// CQL statement for adding message
	addMessage = `INSERT INTO message (mid, chat_id, sender, payload) VALUES (now(), ?, ?, ?)`
	// CQL statement for retrieving message
	retrieveMessages = `SELECT mid, payload, dateOf(mid) as date FROM message WHERE chat_id = ? and sender = ? and id > ?  LIMIT ?`
	// CQL statement for retrieving message
	retrieveOldestMessages = `SELECT mid, payload, dateOf(mid) FROM message WHERE chat_id = ? and sender = ? LIMIT ?`
	// CQL statement for counting message
	countMessages = `SELECT COUNT(*), MAX(mid) FROM message WHERE chat_id = ? and sender = ? and id >= ? and id <= ? ORDER BY id desc`
	// CQL statement for retrieving message in fixed range
	retrieveMessagesByFixedRange = `SELECT mid, payload, dateOf(mid) as date FROM message WHERE chatid = ? and receiver = ? and id >= ? and id <= ?`

	// chat list
	createChat              = `UPDATE chat_list SET chat[?] = ? WHERE owner = ? and type = ?`
	retrieveChatListByType  = `SELECT * FROM chat_list WHERE owner = ? and type = ?`
	retrieveChatLists       = `SELECT * FROM chat_list WHERE owner = ?`
	updateChatLastMessageId = `UPDATE chat_list SET last[?] = ? WHERE owner = ? and type = ?`
	removeChat              = `DELETE chat[?], last[?] FROM chat_list WHERE owner = ? and type = ?`

	// chat meta
	retrieveChatMeta            = `SELECT * FROM chat_meta WHERE chat_id = ? and owner = ?`
	updateChatLastReadMessageId = `INSERT INTO chat_meta (chat_id, owner, last_read) VALUES (?, ?, ?)`
	retrieveLastReadMessageId   = `SELECT last_read, dateOf(last_read) FROM chat_meta WHERE chat_id = ? and owner = ?`
)

//func ListChat(owner string) []string {
//	session := core.NewCassandraRConn()
//	if err := session.Query(listChatId, owner).Exec(); err != nil {
//		result = err
//	}
//}
func AddMessage(msg *Message, payload string, ttl int) (result error) {
	session := core.NewCassandraWConn()
	if err := session.Query(addMessage,
		msg.Receiver, msg.Sender, payload, ttl).Exec(); err != nil {
		result = err
	}
	return
}

// less to more
func GetMessages(receiver string, sender string, last string, limit int) ([]*MessageWithId, error) {
	session := core.NewCassandraRConn()
	iter := session.Query(retrieveMessages,
		receiver, sender, last, limit).Iter()
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

func GetMessagesFromBeginning(receiver string, sender string, limit int) ([]*MessageWithId, error) {
	session := core.NewCassandraRConn()
	iter := session.Query(retrieveOldestMessages,
		receiver, sender, limit).Iter()
	result := make([]*MessageWithId, iter.NumRows(), iter.NumRows())
	var id string
	var payload string
	log.Println(len(result))
	log.Println(iter.NumRows())
	index := 0
	for iter.Scan(&id, &payload) {
		result[index] = &MessageWithId{id, payload}
		index++
		//		log.Println(id, payload)
	}

	return result, iter.Close()
}

func ConfirmMessages(receiver string, sender string, start string, end string, limit int, counts int) ([]*MessageWithId, error) {
	session := core.NewCassandraCConn()
	var dbCounts int
	var endId string
	err := session.Query(countMessages,
		receiver, sender, start, end).Scan(&dbCounts, &endId)

	if err != nil {
		return nil, err
	}
	if dbCounts == counts && endId == end {
		// normal
		//		wsession := core.NewCassandraWConn()
		//		wsession.Query(``)
		log.Println("Confirm messages owner:", receiver, " sender:", sender, " from:", start, " to:", end, "counts:", counts)
		return nil, nil
	} else {
		// unexpected
		iter := session.Query(retrieveMessagesByFixedRange,
			receiver, sender, start, end).Iter()
		result := make([]*MessageWithId, iter.NumRows())
		var id string
		var payload string
		index := 0
		for iter.Scan(&id, &payload) {
			result[index] = &MessageWithId{id, payload}
			index++
			//			log.Println("Confirm ", id, payload)
		}
		if length := len(result); length != 0 && length == counts && result[length-1].Id == end && result[0].Id == start {
			return nil, iter.Close()
		}

		return result, iter.Close()
	}

}

//func GetLatestMessage(receiver string, sender string, limit int) []string {
//	session := core.NewCassandraCConn()
//	iter := session.Query(checkMessage, receiver, sender, limit).Iter()
//	result := make([]string, iter.NumRows())
//	var payload string
//	index := 0
//	for iter.Scan(&payload) {
//		result[index] = payload
//		index++
//	}
//	return result
//}
