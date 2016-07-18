// operation
package control

import (
	"bytes"
	"core"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

var routes = [...]*requestMapping{
	&requestMapping{"/ping", addJWT, false, "GET"},
	&requestMapping{"/ping", addJWT, false, "POST"},
	&requestMapping{"/", serveHome, false, "GET"},
	&requestMapping{"/shutdown", shutdown, false, "GET"},
	&requestMapping{"/ws", serveWs, true, "GET"},
}

type requestMapping struct {
	path       string
	handleFunc func(http.ResponseWriter, *http.Request)
	auth       bool
	method     string
}

func addJWT(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer
	jtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Now().Unix(),
		"exp": time.Now().AddDate(0, 0, 1).Unix(),
	})

	tokenString, err := jtoken.SignedString([]byte(core.NewConfig().JWTConfig.Secret))
	if err != nil {
		log.Println(err.Error())

	}
	buffer.WriteString("bearer ")
	buffer.WriteString(tokenString)

	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{Name: "Authentication", Value: buffer.String(), Expires: expiration}
	http.SetCookie(w, &cookie)
	log.Println("Token", tokenString)
	vars := mux.Vars(r)
	category := vars["ha"]
	buffer.Reset()
	buffer.WriteString(category)
	w.Write(buffer.Bytes())
}

var homeTemplate = template.Must(template.ParseFiles("index.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {

	//	if r.URL.Path != "/" {
	//		http.Error(w, "Not found", 404)
	//		return
	//	}

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
