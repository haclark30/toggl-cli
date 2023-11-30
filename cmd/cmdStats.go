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
	// setup the start time
	startTime := time.Now()
	var prevStartTime time.Time
	var prevEndTime time.Time

	switch timeFrame {
	case Day:
		prevStartTime = startTime.AddDate(0, 0, -1)
		prevEndTime = prevStartTime
	case Week:
		for startTime.Weekday() != time.Monday {
			startTime = startTime.AddDate(0, 0, -1)
		
		}
		prevStartTime = startTime.AddDate(0, 0, -7)
		prevEndTime = prevStartTime.AddDate(0, 0, 6)
	}

	// get summaries from TogglAPI
	projectSummaries, err := client.GetProjectSummary(workspaceId, startTime, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	prevProjectSummaries, err := client.GetProjectSummary(workspaceId, prevStartTime, prevEndTime)
	if err != nil {
		log.Fatalln(err)
	}

	current, err := client.GetCurrentTimeEntry()
	if err != nil {
		log.Fatalln(err)
	}

	currDuration := time.Now().Sub(current.Start).Truncate(time.Second).Seconds()

	// handle case when current timer has no project
	if current.ProjectID == nil {
		current.ProjectID = new(int)
		*current.ProjectID = 0
		projSum := toggl.ProjectSummary{
			ProjectId:      0,
			TrackedSeconds: int(currDuration),
			UserId:         current.UserID,
		}
		projectSummaries = append(projectSummaries, projSum)
		projMap[0] = toggl.Project{Name: "<no project>", Color: "ffffff"}
	} else {
		// check if current timer is part of project already in the summaries
		foundCurrent := false
		for i, entry := range projectSummaries {
			if entry.ProjectId == *current.ProjectID {
				projectSummaries[i].TrackedSeconds += int(currDuration)
				foundCurrent = true
			}
		}

		if !foundCurrent {
			currProj := toggl.ProjectSummary{
				ProjectId:      *current.ProjectID,
				TrackedSeconds: int(currDuration),
				UserId:         current.UserID,
			}
			projectSummaries = append(projectSummaries, currProj)
		}
	}

	// sort by duration
	sort.Slice(projectSummaries, func(i, j int) bool {
		return projectSummaries[i].TrackedSeconds > projectSummaries[j].TrackedSeconds
	})

	sort.Slice(prevProjectSummaries, func(i, j int) bool {
		return prevProjectSummaries[i].TrackedSeconds > prevProjectSummaries[j].TrackedSeconds
	})

	// write the results
	// TODO: change 25 to the max project name length + 1?
	writer := ansiterm.NewTabWriter(os.Stdout, 25, 5, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(writer, prevStartTime.Format("2006-01-02"))
	writeProjectStats(writer, prevProjectSummaries, projMap, nil)
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, startTime.Format("2006-01-02"))
	writeProjectStats(writer, projectSummaries, projMap, &current)
	writer.Flush()
}

// interface that wraps both TabWriter and io.Writer
type Writer interface {
	Flush() error
	Write([]byte) (int, error)
}

// writes the list of projects using the given Writer
func writeProjectStats(writer Writer, projects []toggl.ProjectSummary, projMap map[int]toggl.Project, current *toggl.TimeEntry) {
	var totalTime time.Duration = 0
	for _, entry := range projects {
		dur := time.Duration(entry.TrackedSeconds * int(time.Second))
		totalTime += dur
	}
	for _, entry := range projects {
		duration := time.Duration(entry.TrackedSeconds * int(time.Second))
		durationStr := duration.String()
		percent := float64(entry.TrackedSeconds) / float64(totalTime.Seconds()) * 100
		bars := strings.Repeat("█", int(percent))
		bars = StringRgb(bars, Hex(projMap[entry.ProjectId].Color))
		coloredProjName := StringRgb(projMap[entry.ProjectId].Name, Hex(projMap[entry.ProjectId].Color))
		durationStr = fmt.Sprintf("%-9s /%6.2f %%", durationStr, percent)
		if current != nil && entry.ProjectId == *current.ProjectID {
			durationStr = StringRgb(durationStr, activeColor)
		}
		fmt.Fprintf(writer, "%s\t%s %s\n", coloredProjName, durationStr, bars)
	}

	fmt.Fprintf(writer, "Total Time\t%s\n", totalTime)
}
