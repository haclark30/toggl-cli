package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/haclark30/toggl-cli/toggl"
)

const configDir = "/.toggl/"
const tokenFile = "api_token"

// Look for config directory at ~/.toggl and check if api_token file exists
func getApiToken() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home dir: %v", err)
	}
	configPath := homeDir + configDir
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating config dir: %v", err)
		}
	}
	tokenPath := configPath + tokenFile
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

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "me":
			handleMe(&client)
		case "current":
			handleCurrentTimer(&client)
		default:
			fmt.Println("not a valid command")
		}
	} else {
		fmt.Println("no command")
	}
}

func handleMe(client *toggl.TogglClient) {
	me, err := client.Me()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(me.Fullname)
}

func handleCurrentTimer(client *toggl.TogglClient) {

	te, err := client.GetCurrentTimeEntry()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(te.Description)
}
