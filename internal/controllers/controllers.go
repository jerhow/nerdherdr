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
	"fmt"
	// "net/url"
	"github.com/gorilla/sessions"
	"github.com/jerhow/nerdherdr/internal/addemployee"
	"github.com/jerhow/nerdherdr/internal/db"
	"github.com/jerhow/nerdherdr/internal/login"
	"github.com/jerhow/nerdherdr/internal/util"
	"github.com/jerhow/nerdherdr/internal/welcome"
	"html/template"
	"net/http"
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
	data := PageData{
		BodyTitle: "Welcome!",
		LoginMsg:  "",
		Common:    util.TmplCommon,
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/index-header-inject.html",
		"templates/header-end.html",
		"templates/header.html",
		"templates/footer.html"))
	tmpl.Execute(w, data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var un, pw string
	// var authenticated bool
	// var userId int

	// fmt.Printf("%#v", r)

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

func AddEmployee_POST(w http.ResponseWriter, r *http.Request) {
	var fname, lname, mi, title, dept, team, hireDate string
	var result bool = false

	loggedIn, userId := util.IsLoggedIn(r)
	if !loggedIn {
		// bounce out
	}

	// fmt.Printf("%#v", r)

	fname = r.PostFormValue("fname")
	lname = r.PostFormValue("lname")
	mi = r.PostFormValue("mi")
	title = r.PostFormValue("title")
	dept = r.PostFormValue("dept")
	team = r.PostFormValue("team")
	hireDate = r.PostFormValue("hire_date")

	// obviously must check for empty values, validate, sanity check, etc
	// result = addemployee.Validate(lname, fname, mi, title, dept, team, hireDate)
	result = addemployee.Validate()
	if result {
		fmt.Println("Validate() call completed")
	}

	// attempt to write to DB
	result = addemployee.PostToDb(lname, fname, mi, title, dept, team, hireDate, userId)
	if result {
		fmt.Println("PostToDb() call completed successfully")
		http.Redirect(w, r, "welcome?um=success", 303)
	}
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
		EmpRows               []db.EmpRow
		UserMsg               string
		NewEmpListOrderBy     string
		EmpListOrderByArrow   string
	}
	data := pageData{
		BodyTitle: "Welcome!",
		Common:    util.TmplCommon,
		UserMsg:   "",
	}

	// A message back to the user
	userMsg := r.URL.Query().Get("um")
	if userMsg == "success" {
		data.UserMsg = "Employee added successfully!"
	}

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

		sortBy := r.URL.Query().Get("sb")
		orderBy := r.URL.Query().Get("ob")
		if orderBy == "0" {
			data.NewEmpListOrderBy = "1"
			// data.EmpListOrderByArrow = "&#8593;"
		} else {
			data.NewEmpListOrderBy = "0"
			// data.EmpListOrderByArrow = "&#8595;"
		}

		data.EmpRows = welcome.FetchEmployeeList(userId, sortBy, orderBy)
		// fmt.Printf("Type: %T \n", db.FetchEmployeeList())
		// fmt.Printf("db.FetchEmployeeList() = %#v \n", db.FetchEmployeeList())
		// fmt.Printf("Type: %T \n", data.EmpRows)
		// fmt.Printf("db.FetchEmployeeList() = %#v \n", data.EmpRows)

		tmpl := template.Must(template.ParseFiles(
			"templates/welcome.html",
			"templates/header-end.html",
			"templates/header.html",
			"templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func AddEmployee_GET(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		Common util.TemplateCommon
	}
	data := pageData{
		Common: util.TmplCommon,
	}

	loggedIn, _ := util.IsLoggedIn(r)

	if loggedIn {
		tmpl := template.Must(template.ParseFiles(
			"templates/add-employee.html",
			"templates/header.html",
			"templates/add-employee-header-inject.html",
			"templates/header-end.html",
			"templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}
