package main

import (
	"bufio"
	"fmt"
	"os"
	"piotrek813/word-bo-piwo/client"
	"piotrek813/word-bo-piwo/notification"
	"time"
)

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
	client.Init()

	bearer, err := client.GetAccessToken()

	if err != nil {
		notification.SendError(err)
	}

	const filename = "list.txt"

	fmt.Printf("bearer: %v\n", bearer)
	for {
		time.Sleep(3 * time.Second)
		exams, err := client.GetPracticalExams(bearer)

		if err != nil {
			notification.SendError(err)
		}

		newItems, err := checkForNewItems(filename, exams)
		if err != nil {
			notification.SendError(err)
		}

		if len(newItems) == 0 {
			fmt.Printf("INFO [%v]: No new items found\n", time.Now())

			continue
		}

		fmt.Println("New items added:")
		for _, item := range newItems {
			notification.Send("Hallelujah", "nowy termin: "+item)
		}

		err = saveListToFile(filename, exams)
		if err != nil {
			notification.SendError(err)
		}
	}
}
