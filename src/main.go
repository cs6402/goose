// main
package main

import (
	"control"
	"core"
	"log"
)

func main() {
	core.NewConfig()
	//	config, err := core.NewConfig("config.toml")
	//	log.Println(&config)
	//	check(err)
	//	configLog()
	//	_, err = core.NewCache()
	//	check(err)
	// TODO database
	// TODO AWS SQS
	// TODO Handler
	wa := make(chan bool)
	go control.NewServer(wa)
	//	go func() {
	//		//		log.Println(&config)
	//		time.Sleep(time.Minute * 1)
	//		wa <- true
	//	}()
	log.Println("Server initialzation succeed!")
	<-wa

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
