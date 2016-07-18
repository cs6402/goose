// database
package repository

import (
	"core"
	"encoding/json"
	"flag"
	. "model"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
)

func TestDatabase(t *testing.T) {
	session := setup(t)
	defer session.Close()
	//	testInsert(t, session)
	testReceive(t, session)
	tearDown(t, session)
}

func testInsert(t *testing.T, session *gocql.Session) {
	msg := &Message{
		"Daniel", "Hi bro", "Admin", uuid.NewV4().String(), time.Now().Unix(), 1,
	}
	payload, _ := json.Marshal(msg)
	if err := AddMessage(msg, string(payload), 6000); err != nil {
		t.Log(err.Error())
		t.Error(err)
	}
}
func testReceive(t *testing.T, session *gocql.Session) {
	//	list, err := GetMessagesFromBeginning("Daniel", 50)
	list, err := GetMessages("Daniel", "bf70b710-4c3d-11e6-a75e-538be8a30143", 50)
	if err != nil {
		t.Log(err.Error())
		t.Error("Failed to get", err)
	}
	t.Log(list)
}

func setup(t *testing.T) *gocql.Session {
	flag.Set("config", "../config.toml")
	_, err := core.NewConfig()
	if err != nil {
		t.Fatal(err.Error())
	}
	session := core.NewCassandraConn()
	//create table
	if err := session.Query(`CREATE TABLE IF NOT EXISTS message (owner varchar, id timeuuid, payload varchar, PRIMARY KEY ((owner),id)) WITH CLUSTERING ORDER BY (id DESC)`).Exec(); err != nil {
		t.Log("can not create table message")
		t.Fatal(err.Error())
	}
	return session
}

func tearDown(t *testing.T, session *gocql.Session) {

	defer session.Close()
	//	if err := session.Query(`Drop Table message`).Exec(); err != nil {
	//		t.Log("can not drop Table testdata")
	//		t.Fatal(err.Error())
	//	}

}
