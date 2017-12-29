package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jerhow/nerdherdr/internal/controllers"
	"github.com/jerhow/nerdherdr/internal/db"
	"github.com/jerhow/nerdherdr/internal/login"
	"github.com/jerhow/nerdherdr/internal/util"
	"log"
	"net/http"
	"os"
)

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT") // There's no way to know this ahead of time
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000" // If running locally. In prod, if we don't get a $PORT value from Heroku, we're fucked anyway.
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Index).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/movies", controllers.AllMovies).Methods("GET")
	r.HandleFunc("/movies", controllers.CreateMovie).Methods("POST")
	r.HandleFunc("/movies", controllers.UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies", controllers.DeleteMovie).Methods("DELETE")
	r.HandleFunc("/movies/{id}", controllers.FindMovie).Methods("GET")
	r.HandleFunc("/tmpl1", controllers.Tmpl1).Methods("GET")
	r.HandleFunc("/tmpl2", controllers.Tmpl2).Methods("GET")

	fmt.Println(db.Doit("Jerry"))
	fmt.Println(util.AddTwoInts(1, 2))
	db.Db1()
	db.DbPopulateStruct()
	db.DbSingleRowQuery()

	fmt.Println("=================")
	// fmt.Println(hashIt("wut"))
	pwd := "pass"
	hash, _ := login.HashPwd(pwd)
	fmt.Println("Password:", pwd)
	fmt.Println("Hash:    ", hash)

	// db.WritePwd(hash)

	match := login.CheckPasswordHash(pwd, hash)
	fmt.Println("Match:   ", match)
	fmt.Println("=================")

	if err := http.ListenAndServe(GetPort(), r); err != nil {
		log.Fatal(err)
	}
}
