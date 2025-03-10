package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"piotrek813/word-bo-piwo/consts"
	"piotrek813/word-bo-piwo/notification"
	"regexp"
	"time"
)

func getTitle(body string) (string, error) {
	re := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("No title found")
}

func GetPracticalExams() ([]string, error) {
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
	req.Header.Set("Authorization", "Bearer "+bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	if resp.StatusCode == http.StatusUnauthorized {
		Init()

		return nil, fmt.Errorf("Token has expired, regenerating")
	}

	body, err := io.ReadAll(resp.Body)

	if title, err := getTitle(string(body)); err == nil {
		err = fmt.Errorf("Invalid reponse expected json got html document with title: %v", title)

		newSleep := consts.DEFAULT_SLEEP * 4
		log.Println("INFO: Sleeping going to sleep for: " + newSleep.String())
		time.Sleep(newSleep)

		notification.SendError(err)
		return nil, err
	}

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
