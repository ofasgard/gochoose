package main

import "github.com/ofasgard/gochoose"

import "fmt"
import "time"
import "math/rand"

func main() {
	rand.Seed(time.Now().UnixNano())
	
	//Initialise a DB and connect to it
	err := gochoose.InitDB("example.db")
	if (err != nil) {
		fmt.Println(err)
		return
	}
	db,err := gochoose.OpenDB("example.db")
	if (err != nil) {
		fmt.Println(err)
		return
	}
	
	//Create five stages: 1 start, 3 intermediate, and one ending.
	start_stage := gochoose.NewStartStage()
	start_stage.Body = "<i>Welcome to the beginning of your adventure! You can see an old well, a beautiful oak tree, and a pit filled with spikes.</i>"
	
	middle_alpha := gochoose.NewStage()
	middle_alpha.Body = "<i>You are in a deep, dark well. Why did you do this?</i>"
	
	middle_beta := gochoose.NewStage()
	middle_beta.Body = "<i>You've climbed a tree. There's nowhere to go from here. That was dumb.</i>"
	
	middle_gamma := gochoose.NewStage()
	middle_gamma.Body = "<i>You... enter the pit of spikes. They're very sharp, and you die shortly thereafter. That was... ill-advised?</i>"
	
	end_stage := gochoose.NewStage()
	end_stage.Body = "<i>Maybe a life of adventuring wasn't for you.</i>"
	
	//Add linkages to the stages.
	start_stage.AddLink(middle_alpha, "Jump into the well.")
	start_stage.AddLink(middle_beta, "Climb up the tree.")
	start_stage.AddLink(middle_gamma, "Blunder into the pit of spikes.")
	middle_alpha.AddLink(end_stage, "Well, I guess this is my life now.")
	middle_beta.AddLink(end_stage, "I don't think I can get down without hurting myself.")
	middle_gamma.AddLink(end_stage, "Seriously, why did you do that?")
	end_stage.AddLink(start_stage, "Let's try that again.")
	
	//Save the new stages.
	gochoose.SaveStage(db, start_stage)
	gochoose.SaveStage(db, middle_alpha)
	gochoose.SaveStage(db, middle_beta)
	gochoose.SaveStage(db, middle_gamma)
	gochoose.SaveStage(db, end_stage)
	
	//Create the server using an example template.
	srv,err := gochoose.NewCYOAServer("", 8080, db, "res/example.html")
	if (err != nil) {
		fmt.Println(err)
		return
	}
	srv.Server.ListenAndServe()
}


