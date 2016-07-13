// server
package control

import (
	"bytes"
	"core"
	"flag"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTemplate = template.Must(template.ParseFiles("index.html"))

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

		var buffer bytes.Buffer
		buffer.WriteString(":")
		buffer.WriteString(core.Get().HttpConfig.Port)

		jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte("My Secret"), nil
			},
			Debug:         true,
			SigningMethod: jwt.SigningMethodHS256,
		})

		n := negroni.Classic()
		mux.HandleFunc("/", serveHome)
		mux.HandleFunc("/shutdown", shutdown)
		mux.HandleFunc("/ws", serveWs)
		mux.Handle("/ping", negroni.New(
			negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
			negroni.Wrap(http.HandlerFunc(shutdown)),
		))
		n.UseHandler(mux)
		server = &graceful.Server{
			Timeout: 10 * time.Second,
			Server: &http.Server{
				Addr:    buffer.String(),
				Handler: n,
			},
			BeforeShutdown: func() bool {
				log.Println("bye")
				return true
			},
		}

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	})
}
