package client

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"piotrek813/word-bo-piwo/notification"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

const (
	loginURL = "https://info-car.pl/oauth2/login"
	authURL  = "https://info-car.pl/oauth2/authorize"
)

// Fetches the CSRF token from the login page
func getCSRFToken() (string, error) {
	resp, err := client.Get(loginURL)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Extract CSRF token using regex
	re := regexp.MustCompile(`name="_csrf" value="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("CSRF token not found")
}

func Login(username, password string) error {
	// Get CSRF token
	csrfToken, err := getCSRFToken()

	log.Println("INFO: got csrfToken: " + csrfToken)

	if err != nil {
		return err
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("_csrf", csrfToken)

	// Create request
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", loginURL)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	fmt.Println("INFO login reponse body: " + string(body))

	return nil
}

// Base64 URL encoding (like JavaScript's btoa with URL-safe modifications)
func base64UrlEncode(data []byte) string {
	a := base64.RawURLEncoding.EncodeToString(data)
	fmt.Println("INFO: generated nouce: " + a)
	return a
}

// Generates a random 45-character nonce (PKCE-compliant characters)
func createNonce() (string, error) {
	const length = 45
	const unreserved = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"

	// Generate random bytes
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Map bytes to unreserved characters
	for i := range bytes {
		bytes[i] = unreserved[bytes[i]%byte(len(unreserved))]
	}

	// Encode to Base64 URL format
	return base64UrlEncode(bytes), nil
}

func GetAccessToken() {
	godotenv.Load()

	Login(os.Getenv("INFO_CAR_LOGIN"), os.Getenv("INFO_CAR_PASSWORD"))

	nouce, err := createNonce()

	// Sample parameters (replace with actual values)
	params := url.Values{
		"response_type": {"id_token token"},
		"client_id":     {"client"},
		"state":         {"exampleState"},
		"redirect_uri":  {"https://info-car.pl/new/assets/refresh.html"},
		"scope":         {"openid profile email resource.read"},
		"nonce":         {nouce},
		"prompt":        {"none"},
	}

	// Construct the full URL
	fullURL := authURL + "?" + params.Encode()

	// Make the GET request with the custom no-redirect client
	resp, err := client.Get(fullURL)
	if err != nil {
		notification.SendError(fmt.Errorf("error making request: %w", err))
		return
	}
	defer resp.Body.Close()

	// Extract the Location header
	location := resp.Header.Get("Location")
	if location == "" {
		notification.SendError(fmt.Errorf("Location header not found"))
		return
	}

	location = strings.Replace(location, "refresh.html#", "?", 1)

	log.Println("INFO: access_token response location header: " + location)

	parsedURL, err := url.Parse(location)
	if err != nil {
		notification.SendError(fmt.Errorf("error parsing location URL: %w", err))
		return
	}

	// Extract the access_token query parameter
	accessToken := parsedURL.Query().Get("access_token")
	if accessToken == "" {
		notification.SendError(fmt.Errorf("access_token not found in the URL"))
		return
	}

	log.Println("INFO: bearer: " + accessToken)
	bearer = accessToken
}
