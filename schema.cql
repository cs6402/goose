CREATE TABLE IF NOT EXISTS chat_message (
	chat_id varchar,
	mid timeuuid,
	sender varchar,
	payload varchar,	
	PRIMARY KEY (
        (chat_id, sender),
        mid
    )
) WITH default_time_to_live=604800;

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

CREATE TABLE IF NOT EXISTS broadcast_message (
	owner varchar,
	mid timeuuid,
	// 1: system, 2: loc, 3: business
	type int,
	broadcast_chat_id varchar,
	PRIMARY KEY (
    		(owner), mid
    )
)

CREATE TABLE IF NOT EXISTS broadcast_message_repository (
	mid timeuuid,
	broadcast_chat_id varchar,
	sender_mid varchar,
	payload varchar,
	PRIMARY KEY (
    		(broadcast_chat_id), mid
    )
)

CREATE TABLE IF NOT EXISTS broadcast_meta (
	broadcast_chat_id varchar,
	owner varchar,
	last_read timeuuid,
	PRIMARY KEY (
    		(broadcast_chat_id), owner, last_read
    )
)

CREATE TABLE IF NOT EXISTS callin_message (
	broadcast_chat_id varchar,
	mid timeuuid,
	payload varchar,
	sender varchar,
	PRIMARY KEY (
    		(broadcast_chat_id), mid
    )
)