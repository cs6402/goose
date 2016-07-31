// database
package core

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gocql/gocql"
)

func TestDatabase(t *testing.T) {
	session, cluster := setup(t)
	defer session.Close()
	testInsert(t, session)
	tearDown(t, cluster)
}

func testInsert(t *testing.T, session *gocql.Session) {
	// insert a testdata
	if err := session.Query(`INSERT INTO testdata (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		t.Log("can not create table testdata")
		t.Fatal(err.Error())
	}

	var id gocql.UUID
	var text string

	if err := session.Query(`SELECT id, text FROM testdata WHERE timeline = ? LIMIT 1`,
		"me").Consistency(gocql.One).Scan(&id, &text); err != nil {
		t.Error(err.Error())
	}
	fmt.Println("Tweet:", id, text)

	// list all tweets
	iter := session.Query(`SELECT id, text FROM testdata WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

func setup(t *testing.T) (*gocql.Session, *gocql.ClusterConfig) {
	cluster := gocql.NewCluster()
	cluster.Hosts = []string{"127.0.0.1"}
	cluster.ProtoVersion = 4
	cluster.Keyspace = "test"
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()

	if err != nil {
		t.Log("can not create session")
		t.Fatal(err.Error())
	}

	// create test keyspace NOT Supported
	//	if err := session.Query(`CREATE KEYSPACE tk WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};`).Exec(); err != nil {
	//		t.Log("can not create test keyspace")
	//		t.Fatal(err.Error())
	//	}

	//create table testdata
	if err := session.Query(`CREATE TABLE IF NOT EXISTS testdata (timeline varchar, id timeuuid, text varchar, PRIMARY KEY (timeline,id,text))`).Exec(); err != nil {
		t.Log("can not create table testdata")
		t.Fatal(err.Error())
	}
	return session, cluster
}

func tearDown(t *testing.T, cluster *gocql.ClusterConfig) {

	cluster.Timeout = 6000 * time.Millisecond
	session, err := cluster.CreateSession()

	if err != nil {
		t.Log("can not create tearDown session")
		t.Fatal(err.Error())
	}
	defer session.Close()
	if err := session.Query(`Drop Table testdata`).Exec(); err != nil {
		t.Log("can not drop Table testdata")
		t.Fatal(err.Error())
	}
	//	if err := session.Query(`Drop KEYSPACE test`).Exec(); err != nil {
	//		t.Log("can not drop ks test")
	//		t.Fatal(err.Error())
	//	}
}
