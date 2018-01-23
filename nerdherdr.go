package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jerhow/nerdherdr/internal/controllers"
	"github.com/jerhow/nerdherdr/internal/db"
	"github.com/jerhow/nerdherdr/internal/util"
	"log"
	"net/http"
	"os"
)

// Middleware. Got this from:
// https://github.com/jonahgeorge/force-ssl-heroku (Copyright (c) 2017 Jonah George - MIT License)
// ...but have changed the name, as well as the ENV var it's looking for.
// If I ever make a meaningful and useful improvement to this, I'll submit a PR to Jonah.
//
// Summary from the author:
// Heroku does SSL termination at its load balancer.
// However, the app can tell if the original request was made with HTTP by inspecting headers
// inserted by Heroku. We can use this to redirect to the HTTPS Heroku url.
//
// Caveat from the author:
// It works because Heroku exposes your app through a reverse proxy which is used
// for load-balancing and other things. This reverse proxy does SSL termination and
// forwards to your app which should only accept connections from localhost.
// The middleware detects this situation by inspecting headers inserted by Heroku's reverse proxy;
// since headers can be spoofed, you should not use this middleware anywhere that's not behind
// such a proxy!
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

	util.Setup()
	db.SetUpEnv()

	r := mux.NewRouter()
	r.HandleFunc("/", controllers.Index).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")
	r.HandleFunc("/welcome", controllers.Welcome).Methods("GET")
	r.HandleFunc("/welcome", controllers.Welcome_POST).Methods("POST")
	r.HandleFunc("/add-employee", controllers.AddEmployee_GET).Methods("GET")
	r.HandleFunc("/add-employee", controllers.AddEmployee_POST).Methods("POST")

	if err := http.ListenAndServe(getPort(), forceSslHeroku(r)); err != nil {
		log.Fatal(err)
	}
}
