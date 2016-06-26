// main
package main

import (
	"core"
	"log"
	"os"
	"time"
)

func main() {
	config, err := core.NewConfig(os.Args[1])
	log.Println(&config)
	check(err)
	configLog()
	_, err = core.NewCache()
	check(err)
	// TODO database
	// TODO AWS SQS
	// TODO Handler
	wa := make(chan bool)
	go func() {
		log.Println(&config)
		time.Sleep(time.Second * 5)
		wa <- true
	}()
	log.Println("Server initialzation succeed!")
	<-wa
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
