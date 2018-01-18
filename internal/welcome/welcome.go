// Package welcome gathers any helper functions which are only relevant to the 'welcome' view
package welcome

import (
	"github.com/jerhow/nerdherdr/internal/db"
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
func FetchEmployeeList(userId int) []db.EmpRow {
	return db.EmployeeList(userId)
}
