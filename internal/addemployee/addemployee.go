package addemployee

import (
	"fmt"
)

func Validate(lname, fname, mi, title, dept, team, hireDate string) bool {
	fmt.Println(lname)
	fmt.Println(fname)
	fmt.Println(mi)
	fmt.Println(title)
	fmt.Println(dept)
	fmt.Println(team)
	fmt.Println(hireDate)
	return true
}

func PostToDb(lname, fname, mi, title, dept, team, hireDate string, userId int) bool {
	fmt.Println(lname)
	fmt.Println(fname)
	fmt.Println(mi)
	fmt.Println(title)
	fmt.Println(dept)
	fmt.Println(team)
	fmt.Println(hireDate)
	fmt.Println(userId)
	return true
}
