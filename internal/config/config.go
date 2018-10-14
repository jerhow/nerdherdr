// Package config is the place for application-level configuration.
// Things like feature flags, etc.
// NOTE: Only put things in here that you would feel comfortable pushing
// to the repo. Nothing sensitive (keys, auth stuff, db crap) should
// be here (those should live in the ENV config variables).
package config

const LOCAL_PORT string = "3000"
const DISPLAY_BRANDING bool = false
const SESSION_LENGTH int = 3600 // 1 hour
const SHOW_NAV bool = true

// "Tools for Technical Managers"
// "Guiding the Human Side of Technology"
const PAGE_TITLE string = "Nerdherdr: Guiding the Human Side of Technology"
const MASTHEAD_TAGLINE string = "Tools for Technical Managers"

const STATIC_ASSET_URL_BASE_PROD string = "https://s3.amazonaws.com/nerdherdr/"
const STATIC_ASSET_URL_BASE_STAGE string = "https://s3.amazonaws.com/nerdherdr/"
const STATIC_ASSET_URL_BASE_DEV string = "http://localhost:8080/"
const STATIC_ASSET_URL_BASE_LOCAL string = "http://localhost:8080/"
const STATIC_ASSET_URL_BASE_DEFAULT string = "http://localhost:8080/"

const HTML_ARROW_01_UP string = "&#8593;"
const HTML_ARROW_01_DN string = "&#8595;"
const HTML_ARROW_01_RT string = "&#8594;"
const HTML_ARROW_01_LF string = "&#8592;"

const PATH_ROOT string = "rest"

// w.Header().Set("Content-Type", "application/json")
// w.Header().Set("Content-Language", "en")
// w.Header().Set("Cache-Control", "no-store, no-cache")
// w.Header().Set("Location", "https://www.nerdherdr.com")
const HTTP_RESP_CONTENT_TYPE string = "application/json"
const HTTP_RESP_CONTENT_LANGUAGE string = "en"
const HTTP_RESP_CACHE_CONTROL string = "no-store, no-cache"
const HTTP_RESP_LOCATION string = "https://www.nerdherdr.com"
