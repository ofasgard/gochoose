package main

import "gochoose"

import "fmt"
import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano())
	//test_db()
	test_server()
}

func test_db() {
	//Initialise a DB and connect to it
	err := gochoose.InitDB("example.db")
	fmt.Println(err)
	db,err := gochoose.OpenDB("example.db")
	fmt.Println(err)
	
	//Create a new stage and a new user, and point the user's progress to the stage
	stage := gochoose.NewStage()
	stage.Body = "<i>Here is some italic text.</i>"
	stage.Links["Google"] = "https://www.google.com"
	user := gochoose.NewUser()
	user.Progress = stage.ID

	//Save it all and close the database
	gochoose.SaveUser(db, user)
	gochoose.SaveStage(db, stage)
	saved_id := user.ID
	db.Close()
	
	//Now load and do a lookup for our user
	db,err = gochoose.OpenDB("example.db")
	fmt.Println(err)
	user,err = gochoose.LoadUser(db, saved_id)
	stage_id := user.Progress
	stage,err = gochoose.LoadStage(db, stage_id)
	fmt.Println(stage)
	db.Close()
}

func test_server() {
	//Initialise a DB and connect to it
	err := gochoose.InitDB("example.db")
	fmt.Println(err)
	db,err := gochoose.OpenDB("example.db")
	fmt.Println(err)
	
	srv := gochoose.NewCYOAServer("", 8080, db)
	srv.Server.ListenAndServe()
}
