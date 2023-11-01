package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
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

	current, err := client.GetCurrentTimeEntry()
	if err != nil {
		log.Fatalln(err)
	}

	currDuration := time.Now().Sub(current.Start).Truncate(time.Second).Seconds()

	foundCurrent := false
	for i, entry := range report {
		if entry.ProjectId == *current.ProjectID {
			report[i].TrackedSeconds += int(currDuration)
			foundCurrent = true
		}
	}

	if !foundCurrent {
		currProj := toggl.ProjectSummary{
			ProjectId:      *current.ProjectID,
			TrackedSeconds: int(currDuration),
			UserId:         current.UserID,
		}
		report = append(report, currProj)
	}

	sort.Slice(report, func(i, j int) bool {
		return report[i].TrackedSeconds > report[j].TrackedSeconds
	})

	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.Debug)
	var totalTime time.Duration = 0
	for _, entry := range report {
		duration := time.Duration(entry.TrackedSeconds * int(time.Second))
		totalTime += duration
		fmt.Fprintf(writer, "%s\t%s\n", projMap[entry.ProjectId].Name, duration)
	}

	fmt.Fprintf(writer, "Total Time\t%s\n", totalTime)
	writer.Flush()
}
