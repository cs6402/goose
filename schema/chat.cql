CREATE TABLE IF NOT EXISTS chat_message (
	chat_id varchar,
	mid timeuuid,
	sender varchar,
	payload varchar,	
	sid varchar,
	svr varchar,
	PRIMARY KEY (
        (chat_id, sender),
        mid
    )
) WITH default_time_to_live=604800;
CREATE INDEX chat_message_by_sid ON chat_message(sid);

CREATE TABLE IF NOT EXISTS chat_list (
	owner varchar,
	// 1: chat, 2: secret
	type int,
	// key: chat_id, value: target
	chat map<varchar, varchar>,
	// key: chat_id, value: last receive mid
	last map<varchar, timeuuid>,
	PRIMARY KEY (owner, type)
);

CREATE TABLE IF NOT EXISTS chat_meta (
	chat_id varchar,
	owner varchar,
	last_read timeuuid,
	PRIMARY KEY (
    		(chat_id, owner)
    )
)