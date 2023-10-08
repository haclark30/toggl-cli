package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new timer.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		wid := 7636849
		handleStart(client, description, wid)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
func handleStart(client *toggl.TogglClient, description string, workspaceID int) {
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
	fmt.Printf("started new entry: %s\n", description)
}
