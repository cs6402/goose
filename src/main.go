// main
package main

import (
	"control"
	"core"
	"log"
)

func main() {
	configLog()
	core.NewConfig()
	core.NewCache()
	core.NewCassandraConn()
	// TODO AWS SQS
	shutdownCh := make(chan bool)
	go control.NewServer(shutdownCh)
	log.Println("Server initialzation succeed!")
	<-shutdownCh
	log.Println("Server shutdown!")
}

func configLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}
