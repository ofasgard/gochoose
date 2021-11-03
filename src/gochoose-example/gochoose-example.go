package main

import "gochoose"

import "fmt"
import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano())
	test_server()
}

func test_server() {
	//Initialise a DB and connect to it
	err := gochoose.InitDB("example.db")
	fmt.Println(err)
	db,err := gochoose.OpenDB("example.db")
	fmt.Println(err)
	
	//Create a simple stage for users to land on.
	stage := gochoose.NewStartStage()
	stage.Body = "<i>Here is some italic text.</i>"
	stage.Links["Google"] = "https://www.google.com"
	gochoose.SaveStage(db, stage)
	
	//Create the server using an example template.
	srv,err := gochoose.NewCYOAServer("", 8080, db, "src/gochoose-example/example.html")
	fmt.Println(err)
	if err == nil {
		srv.Server.ListenAndServe()
	}
}
