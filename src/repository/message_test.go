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
	LIMITION      = 50
	RECEIVER_USER = "Daniel"
	SEND_USER1    = "Admin"
	//	SEND_USER2    = "PIG"
	TTL    = 86400
	NUMBER = 100000
)

func TestDatabase(t *testing.T) {
	session := setup(t)
	defer session.Close()
	//	func() {
	//		for i := 0; i < NUMBER; i++ {
	//			testInsert(t, session, SEND_USER1)
	//		}
	//	}()
	//	func() {
	//		for i := 0; i < NUMBER; i++ {
	//			testInsert(t, session, SEND_USER2)
	//		}
	//	}()
	//	list := testFirstReceive(t, session)
	//	testConfirm(t, session, list)
	last := "9538ac52-5723-11e6-a94f-59061160675e"
	allCount := 0
	accessCount := 0
	var list []*MessageWithId
	for {
		if list != nil && len(list) != 0 {
			last = list[len(list)-1].Id
		}
		list = testReceive(t, session, last)
		if list != nil && len(list) != 0 {
			time.Sleep(10 * time.Microsecond)
			testConfirm(t, session, list)
			allCount += len(list)
			accessCount++
			t.Log("looping ", last, ",", len(list))
			log.Println("ResulC:", allCount, " AccessC:", accessCount, " ")
		} else {
			t.Log("No looping ", last, ",", len(list))
		}
		if allCount >= NUMBER {
			break
		} else if allCount >= 955 {
			last = "9538ac52-5723-11e6-a94f-59061160675e"
		}
		// simulate network latency
		time.Sleep(10 * time.Microsecond)
	}
	tearDown(t, session)
}

func testInsert(t *testing.T, session *gocql.Session, sender string) {

	msg := &Message{
		RECEIVER_USER, time.Now().String(), sender, uuid.NewV4().String(), time.Now().Unix(), 1,
	}
	payload, _ := json.Marshal(msg)
	if err := AddMessage(msg, string(payload), TTL); err != nil {
		t.Log(err.Error())
		t.Error(err)
	}
}
func testConfirm(t *testing.T, session *gocql.Session, list []*MessageWithId) {
	length := len(list)
	//	end := list[0].Id
	end := list[length-1].Id
	start := list[0].Id
	//	start := list[length-1].Id
	list, err := ConfirmMessages(RECEIVER_USER, SEND_USER1, start, end, LIMITION, length)
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
	list, err := GetMessagesFromBeginning(RECEIVER_USER, SEND_USER1, LIMITION)
	if err != nil {
		t.Log(err.Error())
		t.Error("Failed to get", err)
	}
	return list
}
func testReceive(t *testing.T, session *gocql.Session, last string) []*MessageWithId {
	//	list, err := GetMessagesFromBeginning("Daniel", 50)
	list, err := GetMessages(RECEIVER_USER, SEND_USER1, last, LIMITION)
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
	/*if err := session.Query(`CREATE TABLE IF NOT EXISTS message (owner varchar, sender varchar, id timeuuid, payload varchar, PRIMARY KEY ((owner, sender),id)) WITH CLUSTERING ORDER BY (id DESC)`).Exec(); err != nil {
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
