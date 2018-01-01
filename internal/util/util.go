package util

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

var SESSION_KEY = SessionKey()
var SESSION_COOKIE = SessionCookie()

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

// Returns the session key from ENV
func SessionKey() string {
	//
	// NOTE: Session key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	//
	var key string
	var varExists bool

	key, varExists = os.LookupEnv("NH_PROD_SESS_KEY")
	if !varExists {
		key, varExists = os.LookupEnv("NH_LOCALDEV_SESS_KEY")
		if !varExists {
			fmt.Println("main.sessionKey: No suitable ENV variable found")
			os.Exit(1)
		}
	}

	return key
}

// Returns the session cookie name from ENV
func SessionCookie() string {
	var key string
	var varExists bool

	key, varExists = os.LookupEnv("NH_PROD_SESS_COOKIE")
	if !varExists {
		key, varExists = os.LookupEnv("NH_LOCALDEV_SESS_COOKIE")
		if !varExists {
			fmt.Println("main.sessionKey: No suitable ENV variable found")
			os.Exit(1)
		}
	}

	return key
}
