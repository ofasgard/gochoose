package gochoose

import bolt "go.etcd.io/bbolt"
import "github.com/google/uuid"

import "fmt"
import "net/url"
import "net/http"
import "html/template"

type CYOAServer struct {
	Server *http.Server
	DB *bolt.DB
	Template *template.Template
}

type CYOAFields struct {
	Body template.HTML
	Links template.HTML
}

func NewCYOAServer(host string, port int, db *bolt.DB, template_path string) (*CYOAServer,error) {
	srv := CYOAServer{}
	srv.Server = &http.Server{}
	srv.Server.Addr = fmt.Sprintf("%s:%d", host, port)
	http.HandleFunc("/", srv.CYOAHandler)
	srv.DB = db
	t,err := template.ParseFiles(template_path)
	srv.Template = t
	return &srv,err
}

func (s *CYOAServer) CYOAHandler(w http.ResponseWriter, r *http.Request) {
	user := CookieHandler(w, r, s.DB)
	switch r.Method {
		case "GET":
			GetHandler(w, r, s.DB, s.Template, user)
		default:
			fmt.Println(user)
			fmt.Fprintf(w, "NOT IMPLEMENTED")
	}
}

func CookieHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB) User {
	//Check if the client sent a "gochoose-sessid" cookie.
	c,err := r.Cookie("gochoose-sessid")
	if err != nil {
		//If not, create a new user for the client and set the cookie.
		user := NewUser()
		SaveUser(db, user)
		cookie := http.Cookie{ Name: "gochoose-sessid", Value: user.ID.String() }
		http.SetCookie(w, &cookie)
		return user
	}
	//Check if the cookie corresponds to a valid UUID.
	user_id,err := uuid.Parse(c.Value)
	if err != nil {
		//If not, create a new user for the client and set the cookie.
		user := NewUser()
		SaveUser(db, user)
		cookie := http.Cookie{ Name: "gochoose-sessid", Value: user.ID.String() }
		http.SetCookie(w, &cookie)
		return user
	}
	//Check if the parsed UUID corresponds to a valid user.
	user,err := LoadUser(db, user_id)
	if err != nil {
		//If not, create a new user for the client and set the cookie.
		user := NewUser()
		SaveUser(db, user)
		cookie := http.Cookie{ Name: "gochoose-sessid", Value: user.ID.String() }
		http.SetCookie(w, &cookie)
		return user
	}
	return user
}

func GetHandler(w http.ResponseWriter, r *http.Request, db *bolt.DB, tp *template.Template, user User) {
	//Check if the user's progress should be updated based on the request URI.
	ProgressHandler(db, user, r.URL.Query())
	user,err := LoadUser(db, user.ID)
	if err != nil {
		fmt.Fprintf(w, "ERROR: COULD NOT UPDATE USER PROGRESS [%s]", user.ID.String())
	}
	//Check if the progress stage associated with this user actually exists.
	stage_id := user.Progress
	stage, err := LoadStage(db, stage_id)
	//If it doesn't, return a generic error message.
	if (err != nil) {
		fmt.Fprintf(w, "ERROR: USER HAS INVALID PROGRESS STAGE [%s]", stage_id.String())
		return
	}
	//If it does, construct the HTML for this stage.
	fields := CYOAFields{}
	fields.Body = template.HTML(stage.Body)
	fields.Links = template.HTML(stage.GenerateLinks())
	tp.Execute(w, fields)
}

func ProgressHandler(db *bolt.DB, user User, params url.Values) {
	//Check if progress parameter is present.
	progress, ok := params["progress"]
	if (!ok) {
		return
	}
	//Check if progress parameter is a valid UUID.
	progress_id,err := uuid.Parse(progress[0])
	if (err != nil) {
		return
	}
	//Check if progress UUID map to a valid stage.
	stage,err := LoadStage(db, progress_id)
	if (err != nil) {
		return
	}
	//Update the user's progress with the new stage UUID and save the user.
	user.Progress = stage.ID
	SaveUser(db, user)
	return
}
