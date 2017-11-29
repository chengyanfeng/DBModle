package util

import "time"

var PUNCTUATION []string = []string{".", ";", ",", "(", ")"}
var Jdbc_proxy_url string = "http://localhost:4567/sql"
var MODE string = ""
const (
	ROW_LIMIT_PREVIEW int = 50
	ROW_LIMIT_MAX     int = 1000
)
const (
	DEFAULT_HTTP_TIMEOUT time.Duration = 60 * time.Second
	CACHE_TTL_DEFAULT                  = 60
)
const (
	FMT_CSV  string = "csv"
	FMT_JSON string = "json"
)
const (
	DbPos string = "dbpos"
	User string = "user"
	Json string = "json"
	News string = "news"
	Total string = "total"
	Cat string = "cat"
	Media string = "media"
	Cheng string ="cheng"
)

const (
	IP_REGEX string = "((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))"
)

const (
	GENERAL_ERR int = 400
)

