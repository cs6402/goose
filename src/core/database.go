// database
package core

import (
	"log"

	"time"

	"github.com/gocql/gocql"
)

var onceDb sync.Once

func NewDbConn() {
	cluster := gocql.NewCluster()
	cluster.Hosts = GetConfig().DatabaseConfig.Hosts
	cluster.ProtoVersion = 4
	cluster.Keyspace = GetConfig().DatabaseConfig.Keyspace

	session, err := cluster.CreateSession()

	if err != nil {
		t.Log("can not create session")
		t.Fatal(err.Error())
	}
}
