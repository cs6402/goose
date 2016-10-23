package repository

import (
	"core"
	"log"
	. "model"
)

const (
	// CQL statement for adding message
	addMessage = `INSERT INTO chat_message (mid, chat_id, sender, sid, payload) VALUES (now(), ?, ?, ?, ?)`
	// CQL statement for retrieving message

	retrieveMessages = `SELECT mid, payload, dateOf(mid) as date FROM chat_message WHERE chat_id = ? and sender = ? and mid > ?  LIMIT ?`
	// CQL statement for retrieving message , using mid > 1381670f-1dd2-11b2-7f7f-7f7f7f7f7f7f
	// retrieveOldestMessages = `SELECT mid, payload, dateOf(mid) FROM chat_message WHERE chat_id = ? and sender = ? LIMIT ?`
	// CQL statement for counting message
	countMessages = `SELECT COUNT(*), MAX(mid) FROM chat_message WHERE chat_id = ? and sender = ? and mid >= ? and mid <= ?`
	// CQL statement for retrieving message in fixed range
	retrieveMessagesByFixedRange = `SELECT mid, payload, dateOf(mid) as date FROM chat_message WHERE chat_id = ? and receiver = ? and mid >= ? and mid <= ?`

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

func AddMessage(msg *Message, payload string) error {
	session := core.NewCassandraWConn()
	return session.Query(addMessage, msg.Receiver, msg.Sender, msg.SenderMessageId, payload).Exec()
}

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

//func GetMessagesFromBeginning(receiver string, sender string, limit int) ([]*MessageWithId, error) {
//	session := core.NewCassandraRConn()
//	iter := session.Query(retrieveOldestMessages,
//		receiver, sender, limit).Iter()
//	result := make([]*MessageWithId, iter.NumRows(), iter.NumRows())
//	var id string
//	var payload string
//	log.Println(len(result))
//	log.Println(iter.NumRows())
//	index := 0
//	for iter.Scan(&id, &payload) {
//		result[index] = &MessageWithId{id, payload}
//		index++
//		//		log.Println(id, payload)
//	}

//	return result, iter.Close()
//}

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

func CreateChat(chatId string, target string, owner string, chatType int) error {
	session := core.NewCassandraWConn()
	return session.Query(createChat, chatId, target, owner, chatType).Exec()
}

func RetrieveChatListByType(owner string, chatType int) {
	session := core.NewCassandraRConn()
	session.Query(retrieveChatListByType, owner, chatType).Attempts()
}

func RetrieveChatList(owner string) {
	session := core.NewCassandraRConn()
	session.Query(retrieveChatLists, owner).Bind()
}

func UpdateChatLastMessageId(chatId string, lastMessageId string, owner string, chatType int) error {
	session := core.NewCassandraWConn()
	return session.Query(updateChatLastMessageId, lastMessageId, owner, chatType).Exec()
}

func DeleteChatList(chatId string, owner string, chatType int) error {
	session := core.NewCassandraWConn()
	return session.Query(removeChat, chatId, chatId, owner, chatType).Exec()
}
