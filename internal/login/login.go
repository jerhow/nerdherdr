package login

import (
	"github.com/jerhow/nerdherdr/internal/util"
	"golang.org/x/crypto/bcrypt"
)

func Pepper() string {
	// NOTE: We're not really going to do this in the real world
	return "MyRandomPepper123"
}

func HashPwd(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	util.ErrChk(err)
	return string(bytes), err
}

func CheckPasswordHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil // 'CompareHashAndPassword' returns nil on success, or an error on failure
}
