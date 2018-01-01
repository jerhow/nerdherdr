package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jerhow/nerdherdr/internal/controllers"
	"github.com/jerhow/nerdherdr/internal/db"
	"log"
	"net/http"
	"os"
)

// Middleware. Got this from:
// https://github.com/jonahgeorge/force-ssl-heroku
// ...but have changed the ENV var it's looking for.
func forceSslHeroku(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("NH_GO_ENV") == "production" {
			if r.Header.Get("x-forwarded-proto") != "https" {
				sslUrl := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslUrl, http.StatusTemporaryRedirect)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Get the Port from the environment so we can run on Heroku
func getPort() string {
	var port = os.Getenv("PORT") // There's no way to know this ahead of time on Heroku
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000" // If running locally. In prod, if we don't get a $PORT value from Heroku, we're fucked anyway.
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {

	db.SetUpEnv()

	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Index).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")
	r.HandleFunc("/welcome", controllers.Welcome).Methods("GET")
	r.HandleFunc("/movies", controllers.AllMovies).Methods("GET")
	r.HandleFunc("/movies", controllers.CreateMovie).Methods("POST")
	r.HandleFunc("/movies", controllers.UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies", controllers.DeleteMovie).Methods("DELETE")
	r.HandleFunc("/movies/{id}", controllers.FindMovie).Methods("GET")
	r.HandleFunc("/tmpl1", controllers.Tmpl1).Methods("GET")
	r.HandleFunc("/tmpl2", controllers.Tmpl2).Methods("GET")

	if err := http.ListenAndServe(getPort(), forceSslHeroku(r)); err != nil {
		log.Fatal(err)
	}
}
