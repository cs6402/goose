// main
package main

import (
	"log"
)

func main() {
	configLog()
	cache, err := NewCache()
	if err != nil {
		log.Println(err)
	} else {
		client, err := cache.instance.Get()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(client.Cmd("GET", "Lock"))
		}

	}
	log.Println("Server initialzation succeed!")

}

func configLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
}
