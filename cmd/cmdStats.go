package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
)

const activeColor = "d4af37"

type TimeFrame int

const (
	Day TimeFrame = iota
	Week
	Month
	Year
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "get stats, defaults to current day",
	Run: func(cmd *cobra.Command, args []string) {
		handleStats(client, workspaceId, Day)
	},
}

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "stats for the current day",
	Run: func(cmd *cobra.Command, args []string) {
		handleStats(client, workspaceId, Day)
	},
}

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "stats for the current week (start from most recent Monday)",
	Run: func(cmd *cobra.Command, args []string) {
		handleStats(client, workspaceId, Week)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.AddCommand(dayCmd)
	statsCmd.AddCommand(weekCmd)
}

func handleStats(client *toggl.TogglClient, workspaceId int, timeFrame TimeFrame) {
	projects, err := client.GetProjects(workspaceId)
	if err != nil {
		log.Fatalln(err)
	}
	projMap := make(map[int]toggl.Project, len(projects))

	for _, p := range projects {
		projMap[p.ID] = p
	}

	startTime := time.Now()
	switch timeFrame {
	case Day:
		break
	case Week:
		for startTime.Weekday() != time.Monday {
			startTime = startTime.AddDate(0, 0, -1)
		}
	}
	report, err := client.GetProjectSummary(workspaceId, startTime, time.Now())
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

	var totalTime time.Duration = 0
	for _, entry := range report {
		dur := time.Duration(entry.TrackedSeconds * int(time.Second))
		totalTime += dur
	}

	writer := ansiterm.NewTabWriter(os.Stdout, 10, 5, 1, ' ', tabwriter.Debug)
	for _, entry := range report {
		duration := time.Duration(entry.TrackedSeconds * int(time.Second))
		durationStr := duration.String()
		percent := float64(entry.TrackedSeconds) / float64(totalTime.Seconds()) * 100
		bars := strings.Repeat("█", int(percent))
		bars = StringRgb(bars, Hex(projMap[entry.ProjectId].Color))
		coloredProjName := StringRgb(projMap[entry.ProjectId].Name, Hex(projMap[entry.ProjectId].Color))
		durationStr = fmt.Sprintf("%-9s /%6.2f %%", durationStr, percent)
		if entry.ProjectId == *current.ProjectID {
			durationStr = StringRgb(durationStr, activeColor)
		}
		fmt.Fprintf(writer, "%s\t%s %s\n", coloredProjName, durationStr, bars)
	}

	fmt.Fprintf(writer, "Total Time\t%s\n", totalTime)
	writer.Flush()
}
