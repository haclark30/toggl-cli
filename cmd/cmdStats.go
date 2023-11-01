package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "get stats",
	Run: func(cmd *cobra.Command, args []string) {
		handleStats(client, 7636849)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func handleStats(client *toggl.TogglClient, workspaceId int) {
	projects, err := client.GetProjects(workspaceId)
	if err != nil {
		log.Fatalln(err)
	}
	projMap := make(map[int]toggl.Project, len(projects))

	for _, p := range projects {
		projMap[p.ID] = p
	}

	report, err := client.GetProjectSummary(workspaceId, time.Now(), time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	var totalTime time.Duration = 0
	for _, entry := range report {
		duration := time.Duration(entry.TrackedSeconds * int(time.Second))
		totalTime += duration
		fmt.Printf("%s - %s\n", projMap[entry.ProjectId].Name, duration)
	}
	fmt.Printf("Total Time - %s\n", totalTime)
}
