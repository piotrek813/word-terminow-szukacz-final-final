package notification

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func SendError(error error) {
	log.Println("ERROR: " + error.Error())
}

func Send(title string, message string, token string) {
	// Base URL
	baseURL := "https://wirepusher.com/send"

	// Query parameters
	params := url.Values{}
	params.Add("id", os.Getenv("INFO_CAR_LOGIN"))
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
