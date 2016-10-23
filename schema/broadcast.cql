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