package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/haclark30/toggl-cli/toggl"
)

const tokenPath = "./.toggl/api_token"

func getApiToken() string {
	if _, err := os.Stat(tokenPath); err == nil {
		token, err := os.ReadFile(tokenPath)
		if err != nil {
			log.Fatalf("Error reading token file: %v", err)
		}
		return string(token)
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Token file not found, input API token:")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()

		if err != nil {
			log.Fatalf("Error reading token: %v", err)
		}

		token := scanner.Text()
		file, err := os.Create(tokenPath)

		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}

		defer file.Close()
		_, err = file.Write([]byte(token))

		if err != nil {
			log.Fatalf("Error writing file: %v", err)
		}

		return token
	} else {
		log.Fatalf("Error checking token path: %v", err)
		return ""
	}
}

func main() {
	apiToken := getApiToken()
	client := toggl.NewTogglClient(apiToken)
	me, err := client.Me()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(me.Fullname)

	te, err := client.GetCurrentTimeEntry()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(te.Description)
}
