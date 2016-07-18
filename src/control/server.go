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

	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"
)

var addr = flag.String("addr", ":8080", "http service address")

var serverOnce sync.Once
var server *graceful.Server
var shutdownCh chan bool

func NewServer(sh chan bool) {
	serverOnce.Do(func() {
		shutdownCh = sh
		flag.Parse()
		go hub.run()
		r := mux.NewRouter()

		var buffer bytes.Buffer
		buffer.WriteString(":")
		buffer.WriteString(core.NewConfig().HttpConfig.Port)

		jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(core.NewConfig().JWTConfig.Secret), nil
			},
			Debug: true,
			Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
				func(r *http.Request) (string, error) {
					c, err := r.Cookie("Authentication")
					if err != nil {
						return "", err
					}
					authHeaderParts := strings.Split(c.Value, " ")
					if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
						return "", fmt.Errorf("Authorization header format must be Bearer {token}")
					}
					return authHeaderParts[1], nil
				}),

			SigningMethod: jwt.SigningMethodHS256,
		})

		n := negroni.Classic()
		for _, v := range routes {
			if v.auth {
				r.Handle(v.path, negroni.New(
					negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
					negroni.Wrap(http.HandlerFunc(v.handleFunc)),
				)).Methods(v.method)
			} else {
				r.HandleFunc(v.path, v.handleFunc).Methods(v.method)
			}
		}

		n.UseHandler(r)
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
