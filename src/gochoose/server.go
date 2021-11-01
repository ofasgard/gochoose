package gochoose

import "fmt"
import "net/http"

func NewCYOAServer(host string, port int) *http.Server {
	srv := &http.Server{}
	srv.Addr = fmt.Sprintf("%s:%d", host, port)
	http.HandleFunc("/", CYOAHandler)
	return srv
}

func CYOAHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the client sent a "gochoose-sessid" cookie.
	c,err := r.Cookie("gochoose-sessid")
	//If they didn't, create a new user for the client.
	if err != nil {
		//todo
	}
	fmt.Fprintf(w, "NOT IMPLEMENTED")	
}
