// Package controllers is the place for route handler functions, which people
// seem to like referring to as controllers.
// These functions are named in the route declarations for the mux router.
package controllers

import (
	"fmt"
	"github.com/jerhow/nerdherdr/internal/login"
	"html/template"
	"net/http"
	"os"
)

func AllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET /movies (All Movies). Not implemented yet. PORT env var: "+os.Getenv("PORT"))
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	fmt.Fprintf(w, "ID: %s", id)
	// fmt.Fprintln(w, "GET /movies/{id} (Find Movie). Not implemented yet, fool!")
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%+v\n", r)
	id := r.FormValue("id")
	name := r.PostFormValue("name")
	fmt.Fprintf(w, "Hello, %s! ID: %s", name, id)
	// fmt.Fprintln(w, "POST /movies (Create Movie). Not implemented yet.")
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PUT /movies (Update Movie). Not implemented yet.")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DELETE /movies (Delete Movie). Not implemented yet.")
}

func Index(w http.ResponseWriter, r *http.Request) {
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
	type PageData struct {
		PageTitle string
		BodyTitle string
		LoginMsg  string
	}

	var un, pw string
	data := PageData{
		PageTitle: "Nerdherdr: Tools for Technical Managers",
		BodyTitle: "Welcome!",
	}

	un = r.PostFormValue("un")
	pw = r.PostFormValue("pw")

	if un == "" || pw == "" {
		// fmt.Println("un is an empty string, no value provided")
		data.LoginMsg = "Invalid login (one or more fields left blank)"
	} else {
		// fmt.Printf("%+v\n", un)
		if login.Authenticate(un, pw) {
			data.LoginMsg = "Valid login!!! :)"
		} else {
			data.LoginMsg = "Invalid login (auth)"
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)

	// driver := "mysql"
	// dsn := dbDsn()
	// db, err := sql.Open(driver, dsn)
	// util.ErrChk(err)
	// defer db.Close()

	// err = db.Ping()
	// util.ErrChk(err)

	// err = db.QueryRow(
	// 	"SELECT l_name, f_initial FROM t_users WHERE id = ?", 1).Scan(&LName, &FInitial)

	// switch {
	// case err == sql.ErrNoRows:
	// 	fmt.Println("No user with that ID")
	// case err != nil:
	// 	log.Fatal(err)
	// default:
	// 	fmt.Printf("\nUSER: %s, %s\n", LName, FInitial)
	// }
}

// ===========================================================
// This stuff below was from a tutorial. Will clean up soon.
//
type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	BodyTitle string
	Todos     []Todo
}

func Tmpl1(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/tmpl1.html"))
	data := TodoPageData{
		PageTitle: "Template 1",
		BodyTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.Execute(w, data)
}

func Tmpl2(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		PageTitle string
		BodyTitle string
	}
	tmpl := template.Must(template.ParseFiles("templates/tmpl2.html"))
	data := PageData{
		PageTitle: "Template 2",
		BodyTitle: "This is the second template",
	}
	tmpl.Execute(w, data)
}
