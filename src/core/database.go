// database
package core

import (
	"log"

	"sync"
	"time"

	"github.com/gocql/gocql"
)

var onceRDb sync.Once
var onceWDb sync.Once
var onceCDb sync.Once
var dbRConn *gocql.Session
var dbWConn *gocql.Session
var dbCConn *gocql.Session

func NewCassandraRConn() *gocql.Session {
	onceRDb.Do(func() {
		cluster := gocql.NewCluster()
		cluster.Hosts = NewConfig().DatabaseConfig.Hosts
		cluster.ProtoVersion = 4
		cluster.Keyspace = NewConfig().DatabaseConfig.Keyspace
		cluster.Timeout = 6000 * time.Millisecond
		cluster.Consistency = gocql.ParseConsistency(NewConfig().DatabaseConfig.Read)
		session, err := cluster.CreateSession()

		if err != nil {
			log.Panicln("Create read session failed.", err)
		}
		dbRConn = session
	})
	return dbRConn
}
func NewCassandraWConn() *gocql.Session {
	onceWDb.Do(func() {
		cluster := gocql.NewCluster()
		cluster.Hosts = NewConfig().DatabaseConfig.Hosts
		cluster.ProtoVersion = 4
		cluster.Keyspace = NewConfig().DatabaseConfig.Keyspace
		cluster.Timeout = 6000 * time.Millisecond
		cluster.Consistency = gocql.ParseConsistency(NewConfig().DatabaseConfig.Write)
		session, err := cluster.CreateSession()

		if err != nil {
			log.Panicln("Create write session failed.", err)
		}
		dbWConn = session
	})
	return dbWConn
}
func NewCassandraCConn() *gocql.Session {
	onceCDb.Do(func() {
		cluster := gocql.NewCluster()
		cluster.Hosts = NewConfig().DatabaseConfig.Hosts
		cluster.ProtoVersion = 4
		cluster.Keyspace = NewConfig().DatabaseConfig.Keyspace
		cluster.Timeout = 6000 * time.Millisecond
		cluster.Consistency = gocql.ParseConsistency(NewConfig().DatabaseConfig.Confirm)
		session, err := cluster.CreateSession()

		if err != nil {
			log.Panicln("Create confirm session failed.", err)
		}
		dbCConn = session
	})
	return dbCConn
}
