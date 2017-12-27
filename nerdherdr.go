package main

import (
	"./mylib" // "github.com/jerhow/nerdherdr/mylib"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
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

func errChk(err error) {
	if err != nil {
		log.Fatal(err) // panic(err.Error)
	}
}

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

func dbDsn() string {
	dbUser := "jerry"
	dbPass := "pass"
	dbHost := "go_mysql_1"
	dbPort := "3306"
	dbName := "nerdherdr01"

	connStr := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	return connStr
}

func db() {
	fmt.Println("sup from db()")
	driver := "mysql"
	dsn := dbDsn()
	db, err := sql.Open(driver, dsn)
	errChk(err)
	defer db.Close()

	err = db.Ping()
	errChk(err)

	stmtIns, err := db.Prepare("INSERT INTO t_users (l_name, f_initial) VALUES (?, ?)")
	errChk(err)
	defer stmtIns.Close()

	_, err2 := stmtIns.Exec("Franklin", "A")
	errChk(err2)
	defer stmtIns.Close()

	// insert, err := db.Query("INSERT INTO t_users (l_name, f_initial) VALUES ('Franklin', 'A')")
	// if err != nil {
	// 	log.Fatal(err) // panic(err.Error)
	// }
	// defer insert.Close()
}

func dbPopulateStruct() {
	fmt.Println("sup from dbPopulateStruct()")

	// Tag... - a very simple struct
	type Tag struct {
		Id       int    `json:"id"`
		Lname    string `json:"l_name"`
		Finitial string `json:"f_initial"`
	}

	driver := "mysql"
	dsn := dbDsn()
	db, err := sql.Open(driver, dsn)
	errChk(err)
	defer db.Close()

	err = db.Ping()
	errChk(err)

	results, err2 := db.Query("SELECT id, l_name, f_initial FROM t_users WHERE id > ?", 6)
	errChk(err2)

	fmt.Println()
	for results.Next() { // for each row, scan the result into the tag object (struct)
		var tag Tag
		err := results.Scan(&tag.Id, &tag.Lname, &tag.Finitial)
		errChk(err)
		fmt.Println(strconv.Itoa(tag.Id) + ": " + tag.Lname + ", " + tag.Finitial)
	}
	fmt.Println()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/movies", AllMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", CreateMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", UpdateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", FindMovieEndPoint).Methods("GET")
	r.HandleFunc("/tmpl1", RenderTmpl1).Methods("GET")
	r.HandleFunc("/tmpl2", RenderTmpl2).Methods("GET")

	fmt.Println(mylib.Hi("Jerry"))
	fmt.Println(mylib.AddTwoInts(1, 2))
	db()
	dbPopulateStruct()

	if err := http.ListenAndServe(GetPort(), r); err != nil {
		log.Fatal(err)
	}
}
