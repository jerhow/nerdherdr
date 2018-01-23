// Package welcome gathers any helper functions which are only relevant to the 'welcome' view
package welcome

import (
	"github.com/jerhow/nerdherdr/internal/db"
	"strconv"
)

// Takes a userId.
// Returns a boolean indicating whether results were found, and the individual values.
// NOTE: Right now this is just a pass-through function from db.FetchUserProfileInfo(),
// however the expectation is that we will have more functionality which will live
// in this function in the future.
func UserProfileInfo(userId int) (bool, string, string, string, string, string) {
	// lname, fname, mi, title, company string := db.FetchUserProfileInfo(userId)
	return db.FetchUserProfileInfo(userId)
}

// Takes a userId, returns a slice of db.EmpRow structs.
// This is another function which seems pointless, but I prefer for organizational reasons
// (I'd rather not call into the 'db' package directly from a controller if possible)
func FetchEmployeeList(userId int, sortBy string, orderBy string) []db.EmpRow {
	return db.EmployeeList(userId, sortBy, orderBy)
}

func DeleteEmployees(userId int, empIds []string) bool {
	var result bool = false
	if len(empIds) > 0 {
		result = db.DeleteEmployees(userId, empIds)
	}
	return result
}

// Takes the raw 'sort' and 'order' string values which the controller got
// off of the query string, sanity-checks the values, and returns string values
// which are appropriate to be used in the actual query.
// If no values, or garbage values, are passed in on the query string,
// we return safe default values.
func ParseEmpListSortAndOrderQsParams(sort string, order string) (string, string) {
	var sortBy string = ""
	var orderBy string = ""
	var found bool = false

	if sort == "" {
		sort = "1"
	}
	if order == "" {
		order = "0"
	}

	sb, _ := strconv.Atoi(sort)
	ob, _ := strconv.Atoi(order)

	sortByMap := map[int]string{
		0: "id",
		1: "lname",
		2: "fname",
		3: "mi",
		4: "title",
		5: "dept",
		6: "team",
		7: "hire_date",
	}
	if sortBy, found = sortByMap[sb]; !found {
		sortBy = "lname"
	}

	if ob == 1 {
		orderBy = "DESC"
	} else { // 0 or default
		orderBy = "ASC"
	}

	return sortBy, orderBy
}
