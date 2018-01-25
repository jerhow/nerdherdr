// Package util is the place for utility functions,
// common variables, constants, data structures, etc
// which are used across multiple parts of the application.
// Note the relationship with the 'config' package, which is
// strictly for hard-coded, common values.
package util

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jerhow/nerdherdr/internal/config"
	"log"
	"net/http"
	"os"
	"time"
)

var STATIC_ASSET_URL_BASE string

var SESSION_KEY = FetchEnvVar("SESS_KEY")
var SESSION_COOKIE = FetchEnvVar("SESS_COOKIE")

type TemplateCommon struct {
	PageTitle          string
	CopyrightYear      int
	StaticAssetUrlBase string
	DisplayBranding    bool
	MastheadTagline    string
}

var TmplCommon TemplateCommon

// Anything we need to initialize in 'util' should go in here,
// or at least be kicked off from in here
func Setup() {

	// Work out the appropriate URL base for static assets
	//
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
		STATIC_ASSET_URL_BASE = config.STATIC_ASSET_URL_BASE_PROD
	case "stage":
		STATIC_ASSET_URL_BASE = config.STATIC_ASSET_URL_BASE_STAGE
	case "dev":
		STATIC_ASSET_URL_BASE = config.STATIC_ASSET_URL_BASE_DEV
	case "devlocal":
		STATIC_ASSET_URL_BASE = config.STATIC_ASSET_URL_BASE_LOCAL
	default:
		STATIC_ASSET_URL_BASE = config.STATIC_ASSET_URL_BASE_DEFAULT
	}

	// Initialize the common template struct
	// Expects STATIC_ASSET_URL_BASE to be set already
	TmplCommon.PageTitle = config.PAGE_TITLE
	TmplCommon.CopyrightYear = time.Now().Year()
	TmplCommon.StaticAssetUrlBase = STATIC_ASSET_URL_BASE
	TmplCommon.DisplayBranding = config.DISPLAY_BRANDING
	TmplCommon.MastheadTagline = config.MASTHEAD_TAGLINE
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

func NoSession(w http.ResponseWriter, r *http.Request) {
	url := "/?um=nosession"
	httpStatusCode := 303
	http.Redirect(w, r, url, httpStatusCode)
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
