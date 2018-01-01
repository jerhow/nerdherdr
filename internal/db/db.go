package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Imports the package solely for its side-effects
	"github.com/jerhow/nerdherdr/internal/util"
	"log"
	// "os"
	"strconv"
)

var DRIVER string
var DB_USER string
var DB_PASS string
var DB_HOST string
var DB_PORT string
var DB_NAME string

// Reads ENV variables from the host environment, and sets up our
// "constants" with the appropriate values for the database and such.
func SetUpEnv() {
	DRIVER = util.FetchEnvVar("DB_DRIVER")
	DB_USER = util.FetchEnvVar("DB_USER")
	DB_PASS = util.FetchEnvVar("DB_PASS")
	DB_HOST = util.FetchEnvVar("DB_HOST")
	DB_PORT = util.FetchEnvVar("DB_PORT")
	DB_NAME = util.FetchEnvVar("DB_NAME")
}

func dsn() string {
	return DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
}

// TEMP
func WritePwd(pwd string) {
	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	stmt, err := dbh.Prepare("UPDATE t_users SET un = ?, pw = ? WHERE id = 1")
	util.ErrChk(err)
	defer stmt.Close()

	_, err2 := stmt.Exec("j@h.com", pwd)
	util.ErrChk(err2)
}

func FetchPwdHash(un string) string {
	var pwdHashFromDb string
	var retVal string = ""

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	err = dbh.QueryRow("SELECT pw FROM t_users WHERE un = ?", un).Scan(&pwdHashFromDb)

	switch {
	case err == sql.ErrNoRows:
		// fmt.Println("No user with that ID")
	case err != nil:
		log.Fatal(err) // Fatal is equivalent to Print() followed by a call to os.Exit(1)
	default:
		// fmt.Printf("\nUSER: %s, %s\n", LName, FInitial)
		// fmt.Println("Something happened")
		retVal = pwdHashFromDb
	}

	return retVal
}

func Db1() {
	fmt.Println("sup from db()")
	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	stmtIns, err := dbh.Prepare("INSERT INTO t_users (l_name, f_initial) VALUES (?, ?)")
	util.ErrChk(err)
	defer stmtIns.Close()

	_, err2 := stmtIns.Exec("Franklin", "A")
	util.ErrChk(err2)

	// insert, err := db.Query("INSERT INTO t_users (l_name, f_initial) VALUES ('Franklin', 'A')")
	// if err != nil {
	// 	log.Fatal(err) // panic(err.Error)
	// }
	// defer insert.Close()
}

func DbPopulateStruct() {
	fmt.Println("sup from dbPopulateStruct()")

	// Tag... - a very simple struct
	type Tag struct {
		Id       int    `json:"id"`
		Lname    string `json:"l_name"`
		Finitial string `json:"f_initial"`
	}

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	// stmt, err := db.Prepare("SELECT id, l_name, f_initial FROM t_users WHERE id > ?")
	// util.ErrChk(err)
	// defer stmt.Close()
	//
	// rows, err := stmt.Query(6)
	// util.ErrChk(err)
	// defer rows.Close()
	//
	// Regarding what's above ^^
	// According to this:
	// http://go-database-sql.org/prepared.html and http://go-database-sql.org/retrieving.html
	// """"
	// Go creates prepared statements for you under the covers.
	// A simple db.Query(sql, param1, param2), for example, works by preparing the sql,
	// then executing it with the parameters and finally closing the statement.
	// """"
	// Which means, the more verbose way we're doing it above is good when you want to
	// explicitly manage the prepared statements you're spawning (as in for heavy reuse, or efficiency),
	// but in other cases, the more concise method below is fine, since it does it all for you.
	// See this as well for more details and concerns over efficiency where it concerns
	// prepared statements. Basically, the preparation, execution and closing of the prepared
	// statement constitute three separate round trips to the database (!).
	// https://www.vividcortex.com/blog/2014/11/19/analyzing-prepared-statement-performance-with-vividcortex/
	//
	//
	// THIS IS PERFECTLY FINE AS WELL:
	rows, err := dbh.Query("SELECT id, l_name, f_initial FROM t_users WHERE id > ?", 6)
	util.ErrChk(err)
	defer rows.Close()

	fmt.Println()
	for rows.Next() { // for each row, scan the result into the tag object (struct)
		var tag Tag
		err := rows.Scan(&tag.Id, &tag.Lname, &tag.Finitial)
		util.ErrChk(err)
		fmt.Println(strconv.Itoa(tag.Id) + ": " + tag.Lname + ", " + tag.Finitial)
	}
	fmt.Println()
}

func DbSingleRowQuery() {
	fmt.Println("sup from dbSingleRowQuery()")
	var LName, FInitial string

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	err = dbh.QueryRow(
		"SELECT l_name, f_initial FROM t_users WHERE id = ?", 1).Scan(&LName, &FInitial)

	switch {
	case err == sql.ErrNoRows:
		fmt.Println("No user with that ID")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("\nUSER: %s, %s\n", LName, FInitial)
	}
}
