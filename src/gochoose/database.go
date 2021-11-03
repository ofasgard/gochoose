package gochoose

import bolt "go.etcd.io/bbolt"
import "github.com/google/uuid"

import "fmt"
import "time"
import "encoding/json"

func OpenDB(path string) (*bolt.DB,error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	return db,err
}

// Initialise a database for the first time.
// Calling on a database that already exists will not wipe it; it will simply do nothing.

func InitDB(path string) error {
	db, err := OpenDB(path)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		//Create buckets to keep track of UUIDs and the contents of each "progress stage".
		tx.CreateBucketIfNotExists([]byte("Users"))
		tx.CreateBucketIfNotExists([]byte("Stages"))
		return nil
	})
	return err
}

// User Struct

type User struct {
	ID uuid.UUID // The unique UUID for this user, sent as a cookie by the browser.
	Progress uuid.UUID // A UUID representing the user's current progress in the game. Acts as a "foreign key" to the Stages bucket.
}

func NewUser() User {
	u := User{}
	u.ID = uuid.New()
	u.Progress,_ = uuid.Parse("00000000-0000-0000-0000-000000000000")
	return u
}

func SaveUser(db *bolt.DB, user User) error {
	err := db.Update(func(tx *bolt.Tx) error {
		id_str := user.ID.String()
		progress_str := user.Progress.String()
	
		b := tx.Bucket([]byte("Users"))
		err := b.Put([]byte(id_str), []byte(progress_str))
		return err
	})
	return err
}

func LoadUser(db *bolt.DB, id uuid.UUID) (User,error) {
	user := User{}
	
	err := db.View(func(tx *bolt.Tx) error {
		id_str := id.String()
	
		b := tx.Bucket([]byte("Users"))
		v := b.Get([]byte(id_str))
		
		if v == nil {
			return fmt.Errorf("No such user: %s", id_str)
		}
		
		user.ID = id
		prog,err := uuid.Parse(string(v))
		user.Progress = prog
		
		return err
	})
	
	return user,err
}

// Stage Struct

type Stage struct {
	ID uuid.UUID // The unique UUID for this stage. The "progress" value of the User struct acts as a foreign key for this identified.
	Body string // The contents of this stage's HTML body.
	Links [][]string // A slice of slices, containing the possible options to progress from this stage. Each element contains ["text", "uuid"].
}

func (s Stage) ToJSON() ([]byte,error) {
	data := struct {
		Body string
		Links [][]string
	} { 
		s.Body,
		s.Links,
	}
	jsondata, err := json.Marshal(data)
	return jsondata,err
}

func (s *Stage) FromJSON(j []byte) error {
	data := struct {
		Body string
		Links [][]string
	} { 
		"",
		nil,
	}
	err := json.Unmarshal(j, &data)
	if (err != nil) {
		return err
	}
	
	s.Body = data.Body
	s.Links = data.Links
	return nil
}

func (s Stage) GenerateLinks() string {
	htmlcontent := ""
	for _, values := range s.Links {
		text := values[0]
		link := values[1]
		htmlcontent += fmt.Sprintf("<a href=\"/?progress=%s\">%s</a><br />\n", link, text)
	}
	return htmlcontent
}

func NewStage() Stage {
	s := Stage{}
	s.ID = uuid.New()
	s.Body = "<b>ERROR - NO CONTENT</b>"
	s.Links = make([][]string, 0)
	return s
}

func NewStartStage() Stage {
	stage := NewStage()
	stage.ID,_ = uuid.Parse("00000000-0000-0000-0000-000000000000")
	return stage
}

func SaveStage(db *bolt.DB, stage Stage) error {
	err := db.Update(func(tx *bolt.Tx) error {
		id_str := stage.ID.String()
		jsondata,err := stage.ToJSON()
		
		if (err != nil) {
			return err
		}
	
		b := tx.Bucket([]byte("Stages"))
		err = b.Put([]byte(id_str), jsondata)
		return err
	})
	return err
}

func LoadStage(db *bolt.DB, id uuid.UUID) (Stage,error) {
	stage := Stage{}
	
	err := db.View(func(tx *bolt.Tx) error {
		id_str := id.String()
	
		b := tx.Bucket([]byte("Stages"))
		v := b.Get([]byte(id_str))
		
		if v == nil {
			return fmt.Errorf("No such stage: %s", id_str)
		}
		
		stage.ID = id
		err := stage.FromJSON(v)
		
		return err
	})
	
	return stage,err
}



