// server
package control

import (
	"bytes"
	"core"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
func addJWT(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "Authentication", Value: `bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAiLCJuYmYiOjE0Njg0MzE1NDUsImV4cCI6MTQ2ODQzNTE0NSwiaWF0IjoxNDY4NDMxNTQ1LCJqdGkiOiJpZDEyMzQ1NiJ9.AkI6DnCkBxDgcCNWxmHjQPTCoIjt4m2a2OxMpUqs3MQ`, Expires: expiration}
	http.SetCookie(w, &cookie)

	//	w.Header().Set("Authentication", `bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2p3dC1pZHAuZXhhbXBsZS5jb20iLCJzdWIiOiJtYWlsdG86bWlrZUBleGFtcGxlLmNvbSIsIm5iZiI6MTQzMDc3OTMwNSwiZXhwIjoxNDYyMzE1MzA1LCJpYXQiOjE0MzA3NzkzMDUsImp0aSI6ImlkMTIzNDU2IiwidHlwIjoiaHR0cHM6Ly9leGFtcGxlLmNvbS9yZWdpc3RlciJ9.KbVlagrOLiy-R65eUrVuno_IAjW-J5i_ySoSrs2SgjU
	//									`)
	var buffer bytes.Buffer
	buffer.WriteString("OK")
	w.Write(buffer.Bytes())
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
				return []byte("616161"), nil
			},
			Debug: true,
			Extractor: func(r *http.Request) (string, error) {
				c, err := r.Cookie("Authentication")
				if err != nil {
					return "", err
				}
				authHeaderParts := strings.Split(c.Value, " ")
				if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
					return "", fmt.Errorf("Authorization header format must be Bearer {token}")
				}
				return authHeaderParts[1], nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})

		n := negroni.Classic()
		mux.HandleFunc("/", serveHome)
		mux.HandleFunc("/shutdown", shutdown)
		mux.HandleFunc("/ws", serveWs)
		mux.HandleFunc("/ping", addJWT)
		mux.Handle("/pong", negroni.New(
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
