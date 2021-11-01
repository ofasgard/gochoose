package gochoose

import bolt "go.etcd.io/bbolt"

import "fmt"
import "net/http"

type CYOAServer struct {
	Server *http.Server
	DB *bolt.DB
}

func NewCYOAServer(host string, port int, db *bolt.DB) *CYOAServer {
	srv := CYOAServer{}
	srv.Server = &http.Server{}
	srv.Server.Addr = fmt.Sprintf("%s:%d", host, port)
	http.HandleFunc("/", srv.CYOAHandler)
	srv.DB = db
	return &srv
}

func (s *CYOAServer) CYOAHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the client sent a "gochoose-sessid" cookie.
	c,err := r.Cookie("gochoose-sessid")
	//If they didn't, create a new user for the client and set the cookie.
	if err != nil {
		user := NewUser()
		SaveUser(s.DB, user)
		cookie := http.Cookie{ Name: "gochoose-sessid", Value: user.ID.String() }
		http.SetCookie(w, &cookie)
	}
	//TODO
	fmt.Fprintf(w, "NOT IMPLEMENTED")	
}
