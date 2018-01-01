package util

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

var SESSION_KEY = FetchEnvVar("SESS_KEY")
var SESSION_COOKIE = FetchEnvVar("SESS_COOKIE")

func ErrChk(err error) {
	if err != nil {
		log.Fatal(err) // panic(err.Error)
	}
}

func IsLoggedIn(r *http.Request) bool {
	// NOTE: Key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	var key = []byte(SESSION_KEY)
	var store = sessions.NewCookieStore(key)
	session, _ := store.Get(r, SESSION_COOKIE)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	} else {
		return true
	}
}

// func LoggedInCheck(w http.ResponseWriter, r *http.Request) {
// 	if !IsLoggedIn(r) {
// 		// fmt.Println("NOT LOGGED IN, REDIRECTING...")
// 		// http.Redirect(w, r, "/login", 403) // 403 Forbidden
// 		http.Error(w, "Forbidden", http.StatusForbidden)
// 	}
// }

// Pass it the ENV variable you want, get back the value.
// This is environment-sensitive (prod, stage, dev, devlocal).
func FetchEnvVar(envVarName string) string {
	var val string
	var varExists bool

	val, varExists = os.LookupEnv("NH_PROD_" + envVarName)
	if !varExists {
		val, varExists = os.LookupEnv("NH_LOCALDEV_" + envVarName)
		if !varExists {
			fmt.Println("util.FetchEnvVar: No suitable ENV variable found for '" + envVarName + "'")
			os.Exit(1)
		}
	}

	return val
}
