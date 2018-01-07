// Package login gathers any helper functions which are only relevant to logging in.
// Currently this package is oriented around bcrypt hashes, but this could change
// in the future.
package welcome

import (
	"github.com/jerhow/nerdherdr/internal/db"
)

func Pepper() string {
	return "Hello"
}

// Takes a userId.
// Returns a boolean indicating whether results were found, and the individual values.
// NOTE: Right now this is just a pass-through function from db.FetchUserProfileInfo(),
// however the expectation is that we will have more functionality which will live
// in this function in the future.
func UserProfileInfo(userId int) (bool, string, string, string, string, string) {
	// lname, fname, mi, title, company string := db.FetchUserProfileInfo(userId)
	return db.FetchUserProfileInfo(userId)
}
