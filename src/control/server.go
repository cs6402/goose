// server
package control

import (
	"bytes"

	"flag"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTemplate = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, r.Host)

}
func shutdown(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/shutdown" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buffer bytes.Buffer
	buffer.WriteString("bye")
	w.Write(buffer.Bytes())
	defer func() {
		server.Stop(10)
		shutdownCh <- true
	}()
}

var serverOnce sync.Once
var server *graceful.Server
var shutdownCh chan bool

func NewServer(sh chan bool) {
	serverOnce.Do(func() {
		shutdownCh = sh
		flag.Parse()
		go hub.run()

		mux := http.NewServeMux()
		mux.HandleFunc("/", serveHome)
		mux.HandleFunc("/shutdown", shutdown)
		mux.HandleFunc("/ws", serveWs)
		//		var buffer bytes.Buffer
		//		buffer.WriteString(":")
		//		buffer.WriteString(core.Config.HttpConfig.Port)

		server = &graceful.Server{
			Timeout: 10 * time.Second,
			Server: &http.Server{
				//				Addr:    buffer.String(),
				Addr:    ":8080",
				Handler: mux,
			},
			BeforeShutdown: func() bool {
				log.Println("bye")
				return true
			},
		}
		server.ListenAndServe()

		//	err := http.ListenAndServe(*addr, nil)
		//	if err != nil {
		//		log.Fatal("ListenAndServe: ", err)
		//	}
	})
}
