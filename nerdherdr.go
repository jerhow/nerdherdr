package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jerhow/nerdherdr/internal/db"
	"github.com/jerhow/nerdherdr/internal/util"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	BodyTitle string
	Todos     []Todo
}

func AllMoviesEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET /movies (All Movies). Not implemented yet. PORT env var: "+os.Getenv("PORT"))
}

func FindMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	fmt.Fprintf(w, "ID: %s", id)
	// fmt.Fprintln(w, "GET /movies/{id} (Find Movie). Not implemented yet, fool!")
}

func CreateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%+v\n", r)
	id := r.FormValue("id")
	name := r.PostFormValue("name")
	fmt.Fprintf(w, "Hello, %s! ID: %s", name, id)
	// fmt.Fprintln(w, "POST /movies (Create Movie). Not implemented yet.")
}

func UpdateMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PUT /movies (Update Movie). Not implemented yet.")
}

func DeleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DELETE /movies (Delete Movie). Not implemented yet.")
}

func RenderTmpl1(w http.ResponseWriter, r *http.Request) {
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

func RenderTmpl2(w http.ResponseWriter, r *http.Request) {
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

func IndexEndPoint(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		PageTitle string
		BodyTitle string
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := PageData{
		PageTitle: "Nerdherdr: Tools for Technical Managers",
		BodyTitle: "Welcome!",
	}
	tmpl.Execute(w, data)
}

func LoginEndPoint(w http.ResponseWriter, r *http.Request) {
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
		data.LoginMsg = "Invalid login"
	} else {
		// fmt.Printf("%+v\n", un)
		data.LoginMsg = "Valid login :)"
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

func pepper() string {
	// NOTE: We're not really going to do this in the real world
	return "MyRandomPepper123"
}

func hashPwd(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	util.ErrChk(err)
	return string(bytes), err
}

func checkPasswordHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil // 'CompareHashAndPassword' returns nil on success, or an error on failure
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexEndPoint).Methods("GET")
	r.HandleFunc("/login", LoginEndPoint).Methods("POST")
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", FindMovieEndPoint).Methods("GET")
	r.HandleFunc("/tmpl1", RenderTmpl1).Methods("GET")
	r.HandleFunc("/tmpl2", RenderTmpl2).Methods("GET")

	// fmt.Println(util.Hi("Jerry"))
	fmt.Println(db.Doit("Jerry"))
	fmt.Println(util.AddTwoInts(1, 2))
	db.Db1()
	db.DbPopulateStruct()
	db.DbSingleRowQuery()

	fmt.Println("=================")
	// fmt.Println(hashIt("wut"))
	pwd := "secret"
	hash, _ := hashPwd(pwd)
	fmt.Println("Password:", pwd)
	fmt.Println("Hash:    ", hash)

	match := checkPasswordHash(pwd, hash)
	fmt.Println("Match:   ", match)
	fmt.Println("=================")

	if err := http.ListenAndServe(GetPort(), r); err != nil {
		log.Fatal(err)
	}
}
