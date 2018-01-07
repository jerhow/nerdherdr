// Package welcome gathers any helper functions which are only relevant to the 'welcome' view.
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
