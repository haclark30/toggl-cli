package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/haclark30/toggl-cli/toggl"
)

const configDir = "/.toggl/"
const tokenFile = "api_token"
const workspaceID = 7636849
const ansiReset = "\033[0m"

type Hex string
type Ansi string

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
		case "start":
			if len(args) > 1 {
				handleStart(&client, args[1])
			} else {
				fmt.Println("start requires additional parameter")
				os.Exit(-1)
			}
		case "rewind":
			if len(args) > 1 {
				handleRewind(&client, args[1])
			} else {
				fmt.Println("rewind requires addtional paramter")
				os.Exit(-1)
			}
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

	proj, err := client.GetProject(te.WorkspaceID, *te.ProjectID)
	if err != nil {
		log.Fatal(err)
	}
	projColor := HextoAnsi(Hex(proj.Color))
	duration := time.Now().Sub(te.Start)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	fmt.Printf(" %s▶ tracking %s on %s for %01d:%02d:%02d%s\n", projColor, te.Description, proj.Name, hours, minutes, seconds, ansiReset)
}

func handleStart(client *toggl.TogglClient, description string) {
	timeEntry := toggl.CreateTimeEntry{
		CreatedWith: "toggl cli",
		WorkspaceID: workspaceID,
		Description: description,
		Start:       time.Now(),
		Duration:    -1,
	}
	err := client.StartTimeEntry(timeEntry)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleRewind(client *toggl.TogglClient, minutes string) {
	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		log.Fatalln(err)
	}

	entries, err := client.GetTimeEntries()
	if err != nil {
		log.Fatalln(err)
	}

	if len(entries) == 0 {
		log.Fatalln("no timers")
	}

	entries[0].Start = entries[0].Start.Add(time.Duration(-minutesInt) * time.Minute)

	// update this to use bulk edit PATCH endpoint instead
	err = client.UpdateTimeEntry(entries[0])

	if err != nil {
		log.Fatalln(err)
	}
	if len(entries) > 1 {
		newStop := entries[1].Stop.Add(time.Duration(-minutesInt) * time.Minute)
		entries[1].Stop = &newStop
		err := client.UpdateTimeEntry(entries[1])
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("updated start time")
}

// The RGB type holds three values: one for red (R), green, (G) and
// blue (B). Each of these colors are on the domain of [0, 255].
type RGB struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
}

// HextoRGB converts a hexadecimal string to RGB values
func HextoRGB(hex Hex) RGB {
	if hex[0:1] == "#" {
		hex = hex[1:]
	}
	r := string(hex)[0:2]
	g := string(hex)[2:4]
	b := string(hex)[4:6]
	R, _ := strconv.ParseInt(r, 16, 0)
	G, _ := strconv.ParseInt(g, 16, 0)
	B, _ := strconv.ParseInt(b, 16, 0)

	return RGB{int(R), int(G), int(B)}
}

// HextoAnsi converts a hexadecimal string to an Ansi escape code
func HextoAnsi(hex Hex) Ansi {
	rgb := HextoRGB(hex)
	str := "\x1b[38;2;" + strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m"
	return Ansi(str)
}
