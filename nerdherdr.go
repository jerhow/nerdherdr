/*
- Create your directory beneath $GOPATH/src
- Create or bring in your source files
- Establish your local repo:
git init
- Make sure you're all good with your Git configs
- Make sure your .gitignore is the way you want it (./vendor/, etc)
- Get a Procfile in place for Heroku. Make your app name the containing folder of the app:
echo "web: your-app-name" > Procfile
- Install your dependencies:
go get -v [...]
- Make sure you have Godeps installed:
go get -uv github.com/tools/godep
- Vendor your dependencies:
godep save ./...
- Verify that it builds and runs as intended
- Make sure you're logged into Heroku:
heroku login
- Make your initial commit, so you're in a clean state
- Create your app in Heroku, with the officially supported Go buildpack:
heroku create -b https://github.com/heroku/heroku-buildpack-go.git
- Push to Heroku for the build and launch:
git push heroku master
- Check the status of your process on Heroku:
heroku ps
- Tail the logs from your process on Heroku:
heroku logs --tail
- Try your request with Curl or Postman or a browser or whatever

=======================================================================================================

Some resources which have been helpful thusfar:
https://devcenter.heroku.com/categories/go-support
https://leanpub.com/howtodeployagowebapptoheroku101/read
https://jonathanmh.com/deploying-golang-app-without-docker/
https://medium.com/@freeformz/hello-world-with-go-heroku-38295332f07b
https://hackernoon.com/build-restful-api-in-go-and-mongodb-5e7f2ec4be94

*/

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET /movies (All Movies). Not implemented yet. PORT env var: "+os.Getenv("PORT"))
}

func FindMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET /movies/{id} (Find Movie). Not implemented yet, fool!")
}

func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "POST /movies (Create Movie). Not implemented yet.")
}

func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PUT /movies (Update Movie). Not implemented yet.")
}

func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DELETE /movies (Delete Movie). Not implemented yet.")
}

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
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", FindMovieEndPoint).Methods("GET")

	if err := http.ListenAndServe(GetPort(), r); err != nil {
		log.Fatal(err)
	}
}
