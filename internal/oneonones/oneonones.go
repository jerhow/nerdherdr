//
// Package oneonones gathers any helper functions which are only relevant to the 'oneonones' view
//
package oneonones

import (
	// "fmt"
	"github.com/jerhow/nerdherdr/internal/db"
)

func FetchEmployeeList(userId int, sortBy string, orderBy string) []db.EmpRow {
	return db.EmployeeList(userId, sortBy, orderBy)
}
