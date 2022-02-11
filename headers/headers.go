package headers

import (
	"net/http"
)

var superProperties string = `eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY2
9yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbi
I6IjEuMC45MDAzIiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Ii
wic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTA0OTY3LC
JjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==`

var userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0"

func Common(req *http.Request) *http.Request {
	req.Header.Set("X-Super-Properties", superProperties)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("accept-language", "en-US")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", userAgent)
	req.Header.Set("TE", "trailers")
	return req
}

func Register(req *http.Request) *http.Request {
	req.Header.Set("accept", "*/*")
	req.Header.Set("authority", "discord.com")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/api/v9/auth/register")
	req.Header.Set("scheme", "https")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("origin", "discord.com")
	req.Header.Set("referer", "discord.com/register")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-Type", "application/json")
	req.Header.Set("user-agent", userAgent)
	req.Header.Set("x-super-properties", superProperties)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	return req
}
