// main
package main

import (
	"control"
	"core"
	"log"
)

func main() {
	core.NewConfig()
	//	log.Println(&config)
	//	check(err)
	//	configLog()
	//	_, err = core.NewCache()
	//	check(err)
	// TODO database
	// TODO AWS SQS
	// TODO Handler
	shutdownCh := make(chan bool)
	go control.NewServer(shutdownCh)
	log.Println("Server initialzation succeed!")
	<-shutdownCh
	log.Println("Server shutdown!")
}
func check(err error) {
	if err != nil {
		log.Println(err)
	} else {
		//		client, err := cache.instance.Get()
		//		if err != nil {
		//			log.Println(err)
		//		} else {
		//			log.Println(client.Cmd("GET", "Lock"))
		//		}
		log.Println("OK")
	}
}
func configLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}
