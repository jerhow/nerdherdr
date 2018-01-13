package util

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
)

var STATIC_ASSET_URL_BASE string

var SESSION_KEY = FetchEnvVar("SESS_KEY")
var SESSION_COOKIE = FetchEnvVar("SESS_COOKIE")

type TemplateCommon struct {
	Value1 string
	Value2 int
}

var TmplCommon = TemplateCommon{
	Value1: "WUT",
	Value2: 42,
}

// Anything we need to initialize in 'util' should go in here,
// or at least be kicked off from in here. I know about Golang's
// package-level 'init' functions, but I want a deliberate level
// of control here, and I want main.main() to kick this off at startup.
func SetUpEnv() {
	// NOTE: The 'localhost' values here refer to the Apache container
	// I am currently using to serve the static assets. It is meant to
	// mimic the production stack, which uses an S3 bucket.
	//
	// The reason we can't just expose a 'static' route and serve
	// a similarly-named directory from the app's file system is
	// that, on Heroku, the dyno and filesystem are ephemeral.
	// If this were sitting on a server or VPS, that would be possible.

	currentEnv := os.Getenv("NH_GO_ENV")
	switch currentEnv {
	case "production":
		STATIC_ASSET_URL_BASE = "https://s3.amazonaws.com/nerdherdr/"
	case "stage":
		STATIC_ASSET_URL_BASE = "https://s3.amazonaws.com/nerdherdr/"
	case "dev":
		STATIC_ASSET_URL_BASE = "http://localhost:8080/"
	case "devlocal":
		STATIC_ASSET_URL_BASE = "http://localhost:8080/"
	default:
		STATIC_ASSET_URL_BASE = "http://localhost:8080/"
	}
}

func ErrChk(err error) {
	if err != nil {
		log.Fatal(err) // panic(err.Error)
	}
}

func IsLoggedIn(r *http.Request) (bool, int) {
	// NOTE: Key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	var key = []byte(SESSION_KEY)
	var store = sessions.NewCookieStore(key)
	session, _ := store.Get(r, SESSION_COOKIE)
	var userId int
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false, -1 // not found, provide bogus userId value
	} else {
		userId = session.Values["userId"].(int)
		return true, userId
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
// Note the order of precedence in environments:
// (look up PROD first, then STAGE (future), then DEV (future) or DEVLOCAL)
// Fails out hard if an appropriate ENV var is not found.
// TODO: Fail more gracefully, and with proper logging.
// ***
// NOTE: This function is specifically meant for the NH_{ENV}_* variables.
// For non-prefixed vars, just fetch them normally with os.Getenv("WHATEVER")
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
