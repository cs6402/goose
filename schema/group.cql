CREATE TABLE IF NOT EXISTS group_list (
	owner varchar,
	// key: group_chat_id, value: target
	chat map<varchar, varchar>,
	// key: group_chat_id, value: last receive mid
	last map<varchar, timeuuid>,
	// key: group_chat_id, value: the most previous receive mid
	first map<varchar, timeuuid>,
	PRIMARY KEY (owner, type)
);

CREATE TABLE IF NOT EXISTS group_message (
	group_chat_id varchar,
	mid timeuuid,
	sender varchar,
	sender_mid varchar,
	payload varchar,
	PRIMARY KEY (
        (group_chat_id),
        mid
    )
);

CREATE INDEX group_message_by_sender ON group_message(sender);

CREATE TABLE IF NOT EXISTS group_meta (
	group_chat_id varchar,
	owner varchar,
	last_read timeuuid,
	PRIMARY KEY (
    		(group_chat_id), owner, last_read
    )
)