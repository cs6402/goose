// database
package core

import (
	"log"

	"sync"
	"time"

	"github.com/gocql/gocql"
)

var onceDb sync.Once
var dbConn *gocql.Session

func NewCassandraConn() *gocql.Session {
	onceDb.Do(func() {
		cluster := gocql.NewCluster()
		cluster.Hosts = NewConfig().DatabaseConfig.Hosts
		cluster.ProtoVersion = 4
		cluster.Keyspace = NewConfig().DatabaseConfig.Keyspace
		cluster.Timeout = 6000 * time.Millisecond
		session, err := cluster.CreateSession()

		if err != nil {
			log.Panicln("Create session failed.", err)
		}
		dbConn = session
	})
	return dbConn
}
