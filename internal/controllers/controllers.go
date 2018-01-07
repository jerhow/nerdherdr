// Package controllers is the place for route handler functions, which people
// seem to like referring to as controllers.
// These functions are named in the route declarations for the mux router.
package controllers

import (
	"github.com/gorilla/sessions"
	"github.com/jerhow/nerdherdr/internal/login"
	"github.com/jerhow/nerdherdr/internal/util"
	"html/template"
	"net/http"
	"time"
)

var SESSION_KEY = util.FetchEnvVar("SESS_KEY")
var SESSION_COOKIE = util.FetchEnvVar("SESS_COOKIE")

func Index(w http.ResponseWriter, r *http.Request) {

	if util.IsLoggedIn(r) {
		http.Redirect(w, r, "welcome", 303)
	}

	type PageData struct {
		PageTitle string
		BodyTitle string
		LoginMsg  string
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := PageData{
		PageTitle: "Nerdherdr: Tools for Technical Managers",
		BodyTitle: "Welcome!",
		LoginMsg:  "",
	}
	tmpl.Execute(w, data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var un, pw string

	var key = []byte(SESSION_KEY)
	var store = sessions.NewCookieStore(key)
	session, _ := store.Get(r, SESSION_COOKIE)
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}

	type pageData struct {
		PageTitle string
		BodyTitle string
		LoginMsg  string
	}

	data := pageData{
		PageTitle: "Nerdherdr: Tools for Technical Managers",
		BodyTitle: "Please log in",
		LoginMsg:  "",
	}

	un = r.PostFormValue("un")
	pw = r.PostFormValue("pw")

	if un == "" || pw == "" {
		// fmt.Println("un is an empty string, no value provided")
		data.LoginMsg = "Invalid login (one or more fields left blank)"
	} else {
		// fmt.Printf("%+v\n", un)
		if login.Authenticate(un, pw) {
			// data.LoginMsg = "Valid login!!! :)"
			// Set user
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, "welcome", 303)
		} else {
			data.LoginMsg = "Invalid login (auth)"
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// NOTE: Key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	var key = []byte(SESSION_KEY)
	var store = sessions.NewCookieStore(key)
	session, _ := store.Get(r, SESSION_COOKIE)

	// Revoke user's authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		PageTitle          string
		BodyTitle          string
		LoggedIn           string
		CopyrightYear      int
		StaticAssetUrlBase string
	}
	data := pageData{
		PageTitle:          "Nerdherdr: Tools for Technical Managers",
		BodyTitle:          "Welcome!",
		CopyrightYear:      time.Now().Year(),
		StaticAssetUrlBase: util.STATIC_ASSET_URL_BASE,
	}

	if util.IsLoggedIn(r) {
		data.LoggedIn = "Yes"
		tmpl := template.Must(template.ParseFiles("templates/welcome.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		PageTitle     string
		BodyTitle     string
		LoginMsg      string
		CopyrightYear int
	}
	data := pageData{
		PageTitle:     "Nerdherdr: Tools for Technical Managers",
		BodyTitle:     "Welcome!",
		LoginMsg:      "No, but that's okay",
		CopyrightYear: time.Now().Year(),
	}

	tmpl := template.Must(template.ParseFiles("templates/test-header-footer.html", "templates/header.html", "templates/footer.html"))
	tmpl.Execute(w, data)
}
