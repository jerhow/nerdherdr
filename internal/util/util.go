package util

import (
	// "fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func ErrChk(err error) {
	if err != nil {
		log.Fatal(err) // panic(err.Error)
	}
}

func Hi(name string) string {
	return "Hello, " + name
}

func AddTwoInts(x int, y int) int {
	return x + y
}

func IsLoggedIn(r *http.Request) bool {
	// NOTE: Key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	var key = []byte("super-secret-key")
	var store = sessions.NewCookieStore(key)
	session, _ := store.Get(r, "cookie-name")
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
