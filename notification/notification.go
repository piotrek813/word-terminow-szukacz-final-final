package notification

import (
	"fmt"
	"net/http"
	"net/url"
)

func SendError(error error) {
	Send("No i jednak nie działa jak powinno", error.Error())
}

func Send(title string, message string) {
	// Base URL
	baseURL := "https://wirepusher.com/send"

	// Query parameters
	params := url.Values{}
	params.Add("id", "djo2mpG6o")
	params.Add("title", title)
	params.Add("message", message)
	// params.Add("type", "monitoring")
	// params.Add("action", "http://example.tld/server")
	// params.Add("image_url", "http://example.tld/image.png")
	// params.Add("message_id", "12")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, _ := http.Get(fullURL)

	defer resp.Body.Close()
}
