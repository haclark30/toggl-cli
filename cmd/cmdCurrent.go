package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Get current running timer",
	Run: func(cmd *cobra.Command, args []string) {
		handleCurrentTimer(client)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
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
	fmt.Printf(" %s▶ tracking %s on %s for %01d:%02d:%02d%s\n", projColor, te.Description, proj.Name, hours, minutes, seconds, AnsiReset)
}
