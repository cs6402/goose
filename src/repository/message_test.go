// database
package repository

import (
	"core"
	"encoding/json"
	"flag"
	"log"
	. "model"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
)

const (
	LIMITION  = 3
	FROM_USER = "Daniel"
	TO_USER   = "Admin"
	TTL       = 86400
)

func TestDatabase(t *testing.T) {
	session := setup(t)
	defer session.Close()
	go func() {
		for i := 0; i < 10000; i++ {
			testInsert(t, session)
		}
	}()
	//	list := testFirstReceive(t, session)
	//	testConfirm(t, session, list)
	last := "648ef970-51b9-11e6-8b66-a90758ebdcfd"
	allCount := 1
	accessCount := 0
	var list []*MessageWithId
	for {
		if list != nil && len(list) != 0 {
			last = list[0].Id
		}
		list = testReceive(t, session, last)
		if list != nil && len(list) != 0 {
			testConfirm(t, session, list)
			allCount += len(list)
			accessCount++
			t.Log("looping ", last, ",", len(list))
			log.Println("ResulC:", allCount, " AccessC:", accessCount, " ")
		} else {
			t.Log("No looping ", last, ",", len(list))
		}
		if allCount == 10001 {
			break
		}
		time.Sleep(10 * time.Microsecond)
	}
	tearDown(t, session)
}

func testInsert(t *testing.T, session *gocql.Session) {

	msg := &Message{
		FROM_USER, time.Now().String(), TO_USER, uuid.NewV4().String(), time.Now().Unix(), 1,
	}
	payload, _ := json.Marshal(msg)
	if err := AddMessage(msg, string(payload), TTL); err != nil {
		t.Log(err.Error())
		t.Error(err)
	}
}
func testConfirm(t *testing.T, session *gocql.Session, list []*MessageWithId) {
	length := len(list)
	end := list[0].Id
	start := list[length-1].Id
	list, err := ConfirmMessages(FROM_USER, start, end, LIMITION, length)
	if err != nil {
		t.Log(length, ",", end, ",", start)
		t.Log(err.Error())
		t.Error("Failed to get", err)
	}
	if list != nil {
		t.Error("Incorrect message start:", start, " end:", end)
	}
}

func testFirstReceive(t *testing.T, session *gocql.Session) []*MessageWithId {
	//	list, err := GetMessagesFromBeginning("Daniel", 50)
	list, err := GetMessagesFromBeginning(FROM_USER, LIMITION)
	if err != nil {
		t.Log(err.Error())
		t.Error("Failed to get", err)
	}
	return list
}
func testReceive(t *testing.T, session *gocql.Session, last string) []*MessageWithId {
	//	list, err := GetMessagesFromBeginning("Daniel", 50)
	list, err := GetMessages(FROM_USER, last, LIMITION)
	if err != nil {
		t.Log(err.Error())
		t.Error("Failed to get", err)
	}
	return list
}

func setup(t *testing.T) *gocql.Session {
	flag.Set("config", "../config.toml")
	core.NewConfig()
	//	if err != nil {
	//		t.Fatal(err.Error())
	//	}
	session := core.NewCassandraWConn()
	//create table
	/*if err := session.Query(`CREATE TABLE IF NOT EXISTS message (owner varchar, id timeuuid, payload varchar, PRIMARY KEY ((owner),id)) WITH CLUSTERING ORDER BY (id DESC)`).Exec(); err != nil {
		t.Log("can not create table message")
		t.Fatal(err.Error())
	}*/
	return session
}

func tearDown(t *testing.T, session *gocql.Session) {

	defer session.Close()
	//	if err := session.Query(`Drop Table message`).Exec(); err != nil {
	//		t.Log("can not drop Table testdata")
	//		t.Fatal(err.Error())
	//	}

}
