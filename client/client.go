package client

import (
	"net/http"
	"net/http/cookiejar"
)

// Global HTTP client with a cookie jar
var client *http.Client

func Init() *http.Client {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return client
}
