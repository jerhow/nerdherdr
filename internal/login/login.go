// Package login gathers any helper functions which are only relevant to logging in.
// Currently this package is oriented around bcrypt hashes, but this could change
// in the future.
package login

import (
	"github.com/jerhow/nerdherdr/internal/db"
	"github.com/jerhow/nerdherdr/internal/util"
	"golang.org/x/crypto/bcrypt"
)

const BCRYPT_COST int = 14

func Pepper() string {
	// NOTE: We're not really going to do this in the real world
	return "MyRandomPepper123"
}

func HashPwd(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), BCRYPT_COST)
	util.ErrChk(err)
	return string(bytes), err
}

func CheckPasswordHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil // 'CompareHashAndPassword' returns nil on success, or an error on failure
}

func Authenticate(un string, pwdFromUser string) bool {
	pwdHashFromDb := db.FetchPwdHash(un)
	return CheckPasswordHash(pwdFromUser, pwdHashFromDb)
}
