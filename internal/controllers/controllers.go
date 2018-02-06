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
	"github.com/jerhow/nerdherdr/internal/config"
	"github.com/jerhow/nerdherdr/internal/db"
	"github.com/jerhow/nerdherdr/internal/employees"
	"github.com/jerhow/nerdherdr/internal/login"
	"github.com/jerhow/nerdherdr/internal/oneonones"
	"github.com/jerhow/nerdherdr/internal/util"
	"github.com/jerhow/nerdherdr/internal/welcome"
	"html/template"
	"net/http"
	"regexp"
	"strings"
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
		UserMsg   template.HTML
		Common    util.TemplateCommon
	}
	data := PageData{
		BodyTitle: "Welcome!",
		LoginMsg:  "",
		UserMsg:   template.HTML(""),
		Common:    util.TmplCommon,
	}

	data.Common.ShowNav = false // Overriding the default value

	userMsg := r.URL.Query().Get("um")
	if userMsg == "nosession" {
		data.UserMsg = template.HTML(`<span id="user_msg_content" 
			style="color: red;">Session expired. Please log in again.</span>`)
	} else if userMsg == "noauth" {
		data.UserMsg = template.HTML(`<span id="user_msg_content" 
			style="color: red;">Invalid login. Please try again.</span>`)
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
		MaxAge:   config.SESSION_LENGTH,
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
			// data.LoginMsg = "Invalid login (auth)"
			http.Redirect(w, r, "/?um=noauth", 303)
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

func Welcome_GET(w http.ResponseWriter, r *http.Request) {
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
		UserMsg               template.HTML
	}
	data := pageData{
		BodyTitle: "Welcome!",
		Common:    util.TmplCommon,
		UserMsg:   template.HTML(""),
	}

	// If you want to do a message to the user, you could follow this pattern:
	// data.UserMsg = template.HTML(`<span id="user_msg_content"
	// 		style="color: green;">Employee added successfully</span>`)

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
			"templates/welcome.html",
			"templates/header.html",
			"templates/welcome-header-inject.html",
			"templates/add-employee-header-inject.html",
			"templates/header-end.html",
			"templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		util.NoSession(w, r)
		// http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func AddEmployee_GET(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		Common    util.TemplateCommon
		SortByQs  string
		OrderByQs string
	}
	data := pageData{
		Common: util.TmplCommon,
	}

	sortByQs := r.URL.Query().Get("sb")
	orderByQs := r.URL.Query().Get("ob")

	// For the hidden form fields only
	data.SortByQs = sortByQs
	data.OrderByQs = orderByQs

	loggedIn, _ := util.IsLoggedIn(r)

	if loggedIn {
		tmpl := template.Must(template.ParseFiles(
			"templates/add-employee.html"))
		tmpl.Execute(w, data)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
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

	// From the hidden form fields, for the query string
	sortByQs := r.PostFormValue("hdn_sb")
	orderByQs := r.PostFormValue("hdn_ob")

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
		// NOTE: These sb and ob values will sort the list by ID DESC,
		// which I think is useful so that the employee you just entered is
		// right at the top of the list when you land back at /welcome
		url := "employees?um=add_success&sb=" + sortByQs + "&ob=" + orderByQs
		http.Redirect(w, r, url, 303)
	}
}

func Employees_GET(w http.ResponseWriter, r *http.Request) {
	type pageData struct {
		BodyTitle              string
		LoggedIn               string
		UserId                 int
		Common                 util.TemplateCommon
		EmpRows                []db.EmpRow
		UserMsg                template.HTML
		EmpListSortBy          string
		NewEmpListOrderBy      string
		EmpListArrow_id        template.HTML
		EmpListArrow_lname     template.HTML
		EmpListArrow_fname     template.HTML
		EmpListArrow_mi        template.HTML
		EmpListArrow_title     template.HTML
		EmpListArrow_dept      template.HTML
		EmpListArrow_team      template.HTML
		EmpListArrow_hire_date template.HTML
		SortByQs               string
		OrderByQs              string
	}
	data := pageData{
		BodyTitle:              "Welcome!",
		Common:                 util.TmplCommon,
		UserMsg:                template.HTML(""),
		EmpListArrow_id:        template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_lname:     template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_fname:     template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_mi:        template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_title:     template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_dept:      template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_team:      template.HTML("&nbsp;&nbsp;"),
		EmpListArrow_hire_date: template.HTML("&nbsp;&nbsp;"),
	}

	// A message back to the user
	// NOTE: We're injecting a <span> from the server-side, so that we can
	// control the presentation a bit more. For example, green text for success,
	// red for error states, etc
	userMsg := r.URL.Query().Get("um")
	if userMsg == "add_success" {
		data.UserMsg = template.HTML(`<span id="user_msg_content" 
			style="color: green;">Employee added successfully</span>`)
	} else if userMsg == "delete_success" {
		data.UserMsg = template.HTML(`<span id="user_msg_content" 
			style="color: green;">Employee(s) deleted successfully</span>`)
	} else if userMsg == "delete_error" {
		data.UserMsg = template.HTML(`<span id="user_msg_content" 
			style="color: red;">Error: Employee(s) may not have been deleted successfully</span>`)
	}

	loggedIn, userId := util.IsLoggedIn(r)

	if loggedIn {
		data.LoggedIn = "Yes"
		data.UserId = userId

		sortByQs := r.URL.Query().Get("sb")
		orderByQs := r.URL.Query().Get("ob")

		// For the hidden form fields only
		data.SortByQs = sortByQs
		data.OrderByQs = orderByQs

		sortBy, orderBy := employees.ParseEmpListSortAndOrderQsParams(sortByQs, orderByQs)
		data.EmpListSortBy = sortBy
		var arrow template.HTML
		if orderByQs == "0" || orderByQs == "" {
			data.NewEmpListOrderBy = "1"
			arrow = template.HTML(config.HTML_ARROW_01_UP)
		} else {
			data.NewEmpListOrderBy = "0"
			arrow = template.HTML(config.HTML_ARROW_01_DN)
		}

		switch sortBy {
		case "id":
			data.EmpListArrow_id = arrow
		case "lname":
			data.EmpListArrow_lname = arrow
		case "fname":
			data.EmpListArrow_fname = arrow
		case "mi":
			data.EmpListArrow_mi = arrow
		case "title":
			data.EmpListArrow_title = arrow
		case "dept":
			data.EmpListArrow_dept = arrow
		case "team":
			data.EmpListArrow_team = arrow
		case "hire_date":
			data.EmpListArrow_hire_date = arrow
		}

		data.EmpRows = employees.FetchEmployeeList(userId, sortBy, orderBy)
		// fmt.Printf("Type: %T \n", db.FetchEmployeeList())
		// fmt.Printf("db.FetchEmployeeList() = %#v \n", db.FetchEmployeeList())
		// fmt.Printf("Type: %T \n", data.EmpRows)
		// fmt.Printf("db.FetchEmployeeList() = %#v \n", data.EmpRows)

		tmpl := template.Must(template.ParseFiles(
			"templates/employees.html",
			"templates/header.html",
			"templates/employees-header-inject.html",
			"templates/add-employee-header-inject.html",
			"templates/header-end.html",
			"templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		util.NoSession(w, r)
		// http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func Employees_POST(w http.ResponseWriter, r *http.Request) {

	loggedIn, userId := util.IsLoggedIn(r)
	if !loggedIn {
		// bounce out
	}

	// Read the POST keys, looking for 'del_' + empId (ex: 'del_42'),
	// which indicate the delete check boxes. We should end up with a
	// slice of emp IDs as strings, which we can then pass along for processing.
	var str string
	empIds := make([]string, 0)
	re := regexp.MustCompile("^del_\\d+$") // Note that you have to escape the escapes
	r.ParseForm()                          // populates r.Form
	for key, _ := range r.Form {
		if re.MatchString(key) {
			str = strings.Split(key, "_")[1]
			empIds = append(empIds, str)
		}
	}

	// From the hidden form fields, for the query string
	sortByQs := r.PostFormValue("hdn_sb")
	orderByQs := r.PostFormValue("hdn_ob")

	// fmt.Println(empIds)
	result := employees.DeleteEmployees(userId, empIds)
	userMsg := ""
	if !result {
		userMsg = "delete_error"
	} else {
		userMsg = "delete_success"
	}

	url := "employees?um=" + userMsg + "&sb=" + sortByQs + "&ob=" + orderByQs

	http.Redirect(w, r, url, 303)
}

func OneOnOnes_GET(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Hi :)")

	type pageData struct {
		Common  util.TemplateCommon
		EmpRows []db.EmpRow
	}
	data := pageData{
		Common: util.TmplCommon,
	}

	// In case you need it
	// someOsParam := r.URL.Query().Get("someOsParam")

	loggedIn, userId := util.IsLoggedIn(r)

	data.EmpRows = oneonones.FetchEmployeeList(userId, "lname", "ASC")
	fmt.Println(data.EmpRows)

	if loggedIn {
		tmpl := template.Must(template.ParseFiles(
			"templates/oneonones.html",
			"templates/header.html",
			"templates/oneonones-header-inject.html",
			"templates/header-end.html",
			"templates/footer.html"))
		tmpl.Execute(w, data)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}
