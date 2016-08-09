// group
package repository

const (
	// TODO
	addOrCreateGroupChat   = `UPDATE group_list SET chat[?] = ?, first[?] = now() WHERE owner = ?`
	updateGroupLastReceive = `UPDATE group_list SET last[?] = ? WHERE owner = ?`
	retrieveGroupList      = `SELECT * FROM group_list WHERE owner = ?`
	removeGroupChat        = `DELETE chat[?], last[?], first[?] FROM group_list WHERE owner = ?`

	addGroupMessage               = `INSERT INTO group_message (mid, group_chat_id, sender, sender_mid, payload) VALUES (now(), ? , ?, ?, ?)`
	retrieveGroupMessages         = `SELECT * FROM group_message WHERE group_chat_id = ? and mid > ? LIMIT 50`
	retrieveGroupMessagesBySender = `SELECT * FROM group_message WHERE group_chat_id = ? and sender = ? and mid > ? LIMIT 5`

	updateGroupLastRead = `INSERT INTO group_meta (group_chat_id, owner, last_read) values (?, ?, ?)`
	retrieveGroupMeta   = `SELECT * FROM group_meta WHERE group_chat_id = ?`
	removeGroupMeta     = `DELETE FROM group_meta WHERE group_chat_id = ?`
)
