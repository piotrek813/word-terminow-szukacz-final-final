package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const loginURL = "https://info-car.pl/oauth2/login"

// Global HTTP client with a cookie jar
var client *http.Client

func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

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

func login(username, password string) error {
	// Get CSRF token
	csrfToken, err := getCSRFToken()
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
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

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

// Function to get the Location header from the redirect response
func getAccessToken() (string, error) {
	nouce, err := createNonce()

	// Authorization URL
	authURL := "https://info-car.pl/oauth2/authorize"

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
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Extract the Location header
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("Location header not found")
	}

	location = strings.Replace(location, "refresh.html#", "?", 1)

	parsedURL, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf("error parsing location URL: %w", err)
	}

	// Extract the access_token query parameter
	accessToken := parsedURL.Query().Get("access_token")
	if accessToken == "" {
		return "", fmt.Errorf("access_token not found in the URL")
	}

	return accessToken, nil
}

type Exam struct {
	ID             string `json:"id"`
	Places         int    `json:"places"`
	Date           string `json:"date"`
	Amount         int    `json:"amount"`
	AdditionalInfo string `json:"additionalInfo"`
}

type ScheduledHour struct {
	Time          string `json:"time"`
	TheoryExams   []Exam `json:"theoryExams"`
	PracticeExams []Exam `json:"practiceExams"`
	LinkedExams   []Exam `json:"linkedExamsDto"`
}

type ScheduledDay struct {
	Day            string          `json:"day"`
	ScheduledHours []ScheduledHour `json:"scheduledHours"`
}

type Schedule struct {
	ScheduledDays []ScheduledDay `json:"scheduledDays"`
}

type Reservation struct {
	OrganizationID                 string   `json:"organizationId"`
	IsOskVehicleReservationEnabled bool     `json:"isOskVehicleReservationEnabled"`
	IsRescheduleReservation        bool     `json:"isRescheduleReservation"`
	Category                       string   `json:"category"`
	Schedule                       Schedule `json:"schedule"`
}

func mightIGetSomeBeaerTokenWithThatMate() {
	http.Get("")
}

func getPracticalExams(token string) ([]string, error) {
	fmt.Printf("INFO [%v]: Fetching exam\n", time.Now())

	url := "https://info-car.pl/api/word/word-centers/exam-schedule"

	requestBody := map[string]string{
		"category":  "B",
		"endDate":   "2025-05-12T10:18:32.240Z",
		"startDate": "2025-03-11T11:18:32.240Z",
		"wordId":    "25",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)

	var reservation Reservation
	if err := json.Unmarshal(body, &reservation); err != nil {
		fmt.Println("Error unmarshaling JSON debouncing:", err)

		return nil, err
	}

	practicalExams := []string{}

	for _, day := range reservation.Schedule.ScheduledDays {
		for _, hour := range day.ScheduledHours {
			for _, exam := range hour.PracticeExams {
				practicalExams = append(practicalExams, exam.Date)
			}
		}
	}

	return practicalExams, nil
}

func sendNotification(title string, message string) {

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

// Function to save a list of strings to a file
func saveListToFile(filename string, list []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range list {
		_, err := writer.WriteString(item + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush() // Make sure all data is written to the file
	return nil
}

// Function to read the list from a file
func readListFromFile(filename string) ([]string, error) {
	var list []string
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty list
			return list, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// Function to check if there are any new items in the list compared to the file
func checkForNewItems(filename string, newList []string) ([]string, error) {
	existingList, err := readListFromFile(filename)
	if err != nil {
		return nil, err
	}

	// Convert the existing list into a map for fast lookups
	existingMap := make(map[string]bool)
	for _, item := range existingList {
		existingMap[item] = true
	}

	var newItems []string
	for _, newItem := range newList {
		if !existingMap[newItem] {
			newItems = append(newItems, newItem)
		}
	}

	return newItems, nil
}

func main() {
	login("piotreksmolinski04@gmail.com", "LK`vB}|+Vq\"owX#2,*=tOva&;&QF+M(nH+r3")
	bearer, err := getAccessToken()

	if err != nil {
		sendNotification("Error", err.Error())
	}

	const filename = "list.txt"

	fmt.Printf("bearer: %v\n", bearer)
	for {
		time.Sleep(3 * time.Second)
		exams, err := getPracticalExams(bearer)

		if err != nil {
			sendNotification("Error", err.Error())
		}

		newItems, err := checkForNewItems(filename, exams)
		if err != nil {
			sendNotification("Error", err.Error())
		}

		if len(newItems) == 0 {
			fmt.Printf("INFO [%v]: No new items found\n", time.Now())

			continue
		}

		fmt.Println("New items added:")
		for _, item := range newItems {
			sendNotification("Hallelujah", "nowy termin: "+item)
		}

		err = saveListToFile(filename, exams)
		if err != nil {
			sendNotification("Error", err.Error())
		}

	}
}
