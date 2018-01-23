package addemployee

import (
	"fmt"
	"github.com/jerhow/nerdherdr/internal/db"
)

// func Validate(lname, fname, mi, title, dept, team, hireDate string) bool {
func Validate() bool {
	fmt.Println("Hello from Validate(). This is a stub.")
	return true
}

func PostToDb(lname, fname, mi, title, dept, team, hireDate string, userId int) bool {
	return db.AddEmployee(lname, fname, mi, title, dept, team, hireDate, userId)
}
