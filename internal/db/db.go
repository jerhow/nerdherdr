package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Imports the package solely for its side-effects
	"github.com/jerhow/nerdherdr/internal/util"
	"log"
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

func FetchPwdHashAndUserId(un string) (string, int) {
	var pwdHashFromDb string
	var idFromDb int
	var retHash string = ""
	var retId int = -1

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	err = dbh.QueryRow("SELECT id, pw FROM t_users WHERE un = ?", un).Scan(&idFromDb, &pwdHashFromDb)

	switch {
	case err == sql.ErrNoRows:
		// fmt.Println("No user with that ID")
	case err != nil:
		log.Fatal(err) // Fatal is equivalent to Print() followed by a call to os.Exit(1)
	default:
		// fmt.Printf("\nUSER: %s, %s\n", LName, FInitial)
		// fmt.Println("Something happened")
		retHash = pwdHashFromDb
		retId = idFromDb
	}

	return retHash, retId
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

// We use this in EmployeeList()
type EmpRow struct {
	Id        int
	Lname     string
	Fname     string
	Mi        string
	Title     string
	Dept      string
	Team      string
	Hire_date string
}

// Expects the user's id, and sortBy and orderBy values
// Returns a slice of db.EmpRow structs
// NOTE: Assumes that sortBy and orderBy have been sanity-checked
func EmployeeList(userId int, sortBy string, orderBy string) []EmpRow {

	EmpRows := make([]EmpRow, 0)

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	sql := `
	SELECT 
		id,
		lname,
		fname,
		mi,
		title,
		dept,
		team,
		date_format(hire_date, '%c/%e/%Y') as date_of_hire
	FROM 
		t_employees
	WHERE
		userId = ?
	ORDER BY ` + sortBy + ` ` + orderBy + `;`

	// fmt.Println(sql)

	rows, err := dbh.Query(sql, userId)
	util.ErrChk(err)
	defer rows.Close()

	for rows.Next() { // for each row, scan the result into the EmpRow struct
		var row EmpRow
		err := rows.Scan(&row.Id, &row.Lname, &row.Fname, &row.Mi,
			&row.Title, &row.Dept, &row.Team, &row.Hire_date)
		util.ErrChk(err)
		// then append the struct to the slice, which we will pass into the template
		EmpRows = append(EmpRows, row)
	}

	return EmpRows
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

func DeleteEmployees(userId int, empIds []string) bool {

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	// Build the IN clause manually
	// NOTE: At face value, this may seem insecure, however,
	// the 'Welcome_POST' controller ensures that we only
	// get valid emp IDs in here, and the userId comes straight
	// from the user session. Nothing else ever gets into this function
	// via the controller, and nothing else calls this function.
	// The reason for this manual query building is that I don't
	// see a way that we can pass an array of values to stmtIns.Exec()
	// and have them map to placeholder parameters.
	// See: https://golang.org/pkg/database/sql/#DB.Exec
	// If I find a better way to do this type of batch update (I'd rather
	// not send them as individual DELETE queries), I'll revisit this.
	lastEmpIdIdx := len(empIds) - 1
	inClauseValues := ""
	for idx, empId := range empIds {
		inClauseValues = inClauseValues + empId
		if idx < lastEmpIdIdx {
			inClauseValues = inClauseValues + ", "
		}
	}

	sql := "DELETE FROM t_employees WHERE userId = ? AND id IN (" + inClauseValues + ");"
	// fmt.Println(sql)

	stmtIns, err := dbh.Prepare(sql)
	util.ErrChk(err)
	defer stmtIns.Close()

	_, err2 := stmtIns.Exec(userId)
	util.ErrChk(err2)

	return true
}

// Takes the relevant values for the INSERT
// Returns a boolean indicating success|failure
func AddEmployee(lname, fname, mi, title, dept, team, hire_date string, userId int) bool {

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	sql := `
		INSERT INTO t_employees 
		(userId, lname, fname, mi, title, dept, team, hire_date)
		VALUES 
		(?, ?, ?, ?, ?, ?, ?, ?);`

	stmtIns, err := dbh.Prepare(sql)
	util.ErrChk(err)
	defer stmtIns.Close()

	_, err2 := stmtIns.Exec(userId, lname, fname, mi, title, dept, team, hire_date)
	util.ErrChk(err2)

	return true
}

// Takes a userId
// Returns a boolean indicating whether results were found, and the individual values
func FetchUserProfileInfo(userId int) (bool, string, string, string, string, string) {

	var matchFound bool = true
	var lname, fname, mi, title, company string

	dbh, err := sql.Open(DRIVER, dsn())
	util.ErrChk(err)
	defer dbh.Close()

	err = dbh.Ping()
	util.ErrChk(err)

	query := `
		SELECT lname, fname, mi, title, company
		FROM t_users AS u 
			INNER JOIN t_user_profile AS UP
				ON u.user_profile_id = UP.id
		WHERE
			u.id = ?;`

	err = dbh.QueryRow(query, userId).Scan(&lname, &fname, &mi, &title, &company)
	util.ErrChk(err)

	if err == sql.ErrNoRows {
		matchFound = false
		fmt.Println("No record with that ID")
	}

	return matchFound, lname, fname, mi, title, company
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
