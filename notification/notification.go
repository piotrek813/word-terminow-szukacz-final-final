package notification

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	TOKEN_PIOTREK = "djo2mpG6o"
	TOKEN_AGATA   = "RX75mpG6L"
)

func SendError(error error) {
	log.Println("ERROR: " + error.Error())

	Send("No i jednak nie działa jak powinno", error.Error(), TOKEN_PIOTREK)
}

func Send(title string, message string, token string) {
	// Base URL
	baseURL := "https://wirepusher.com/send"

	// Query parameters
	params := url.Values{}
	params.Add("id", token)
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
