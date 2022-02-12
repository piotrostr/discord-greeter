package headers

import (
	"net/http"
)

// TODO make sure those two match sockets headers
var SuperProperties string = `eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY2
9yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbi
I6IjEuMC45MDAzIiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Ii
wic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTA0OTY3LC
JjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==`

var UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0"

func Common(req *http.Request) *http.Request {
	req.Header.Set("X-Super-Properties", SuperProperties)
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("TE", "trailers")
	return req
}

func Register(req *http.Request) *http.Request {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authority", "discord.com")
	req.Header.Set("Method", "POST")
	req.Header.Set("Path", "/api/v9/auth/register")
	req.Header.Set("Scheme", "https")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("Origin", "discord.com")
	req.Header.Set("Referer", "discord.com/register")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("Accept-language", "en-US,en;q=0.9")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Super-Properties", SuperProperties)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	return req
}

func Invite(req *http.Request) *http.Request {
	req.Header.Set("Accept-Language", "en-US,en-IN;q=0.9,zh-Hans-CN;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://discord.com")
	req.Header.Set("Referer", "https://discord.com/channels/@me")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("X-Super-Properties", SuperProperties)
	return req
}
