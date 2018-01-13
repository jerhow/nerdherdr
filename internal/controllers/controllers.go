// Package controllers is the place for route handler functions, which people
// seem to like referring to as controllers.
//
// My philosophy so far in this project is that these are more like 'view controllers'.
// Roughly, they are meant to:
// 1. Field the inbound request
// 2. Establish the structure of the data for the template
// 3. Call out to a same-named 'helper' package for any additional work that needs doing,
//    in order to prepare the result
// 4. Handle any edge cases or non-standard responses
// 5. Ultimately respond; either as an HTTP response, or by invoking one or more templates to be rendered.
// 6. Anything else I am overlooking, or am not yet aware of.
//
// Anyway, these functions are named/assigned in the route declarations for the mux router, in main.main()
//
package controllers

import (
	// "fmt"
	"github.com/gorilla/sessions"
	"github.com/jerhow/nerdherdr/internal/config"
	"github.com/jerhow/nerdherdr/internal/login"
	"github.com/jerhow/nerdherdr/internal/util"
	"github.com/jerhow/nerdherdr/internal/welcome"
	"html/template"
	"net/http"
	"time"
)

var SESSION_KEY = util.FetchEnvVar("SESS_KEY")
var SESSION_COOKIE = util.FetchEnvVar("SESS_COOKIE")

func Index(w http.ResponseWriter, r *http.Request) {

	loggedIn, _ := util.IsLoggedIn(r)
	if loggedIn {
		http.Redirect(w, r, "welcome", 303)
	}

	type PageData struct {
		BodyTitle string
		LoginMsg  string
		Common    util.TemplateCommon
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := PageData{
		BodyTitle: "Welcome!",
		LoginMsg:  "",
		Common:    util.TmplCommon,
	}
	tmpl.Execute(w, data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var un, pw string
	// var authenticated bool
	// var userId int

	var key = []byte(SESSION_KEY)
	var store = sessions.NewCookieStore(key)
	session, _ := store.Get(r, SESSION_COOKIE)
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}

	type pageData struct {
		BodyTitle string
		LoginMsg  string
		Common    util.TemplateCommon
	}

	data := pageData{
		BodyTitle: "Please log in",
		LoginMsg:  "",
		Common:    util.TmplCommon,
	}

	un = r.PostFormValue("un")
	pw = r.PostFormValue("pw")

	if un == "" || pw == "" {
		// fmt.Println("un is an empty string, no value provided")
		data.LoginMsg = "Invalid login (one or more fields left blank)"
	} else {
		// fmt.Printf("%+v\n", un)
		authenticated, userId := login.Authenticate(un, pw)
		if authenticated {
			// data.LoginMsg = "Valid login!!! :)"
			// Set user
			session.Values["authenticated"] = true
			session.Values["userId"] = userId
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
		BodyTitle             string
		LoggedIn              string
		UserId                int
		UserProfileMatchFound bool
		Lname                 string
		Fname                 string
		MI                    string
		Title                 string
		Company               string
		Common                util.TemplateCommon
	}
	data := pageData{
		BodyTitle: "Welcome!",
		Common:    util.TmplCommon,
	}

	// fmt.Printf("%+v\n", data)

	loggedIn, userId := util.IsLoggedIn(r)

	// Fetch user-specific info for user profile
	if loggedIn {
		matchFound, lname, fname, mi, title, company := welcome.UserProfileInfo(userId)
		if matchFound {
			data.UserProfileMatchFound = true
			data.Lname = lname
			data.Fname = fname
			data.MI = mi
			data.Title = title
			data.Company = company
		} else {
			data.UserProfileMatchFound = false
		}
	}

	if loggedIn {
		data.LoggedIn = "Yes"
		data.UserId = userId
		tmpl := template.Must(template.ParseFiles(
			"templates/welcome.html", "templates/header.html", "templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func AddEmployee(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		PageTitle          string
		CopyrightYear      int
		StaticAssetUrlBase string
		DisplayBranding    bool
	}
	data := pageData{
		PageTitle:          "Nerdherdr: Tools for Technical Managers",
		CopyrightYear:      time.Now().Year(),
		StaticAssetUrlBase: util.STATIC_ASSET_URL_BASE,
		DisplayBranding:    config.DISPLAY_BRANDING,
	}

	loggedIn, _ := util.IsLoggedIn(r)

	if loggedIn {
		tmpl := template.Must(template.ParseFiles(
			"templates/add-employee.html", "templates/header.html", "templates/footer.html"))
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
