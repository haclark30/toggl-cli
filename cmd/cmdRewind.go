package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/spf13/cobra"
)

var rewindCmd = &cobra.Command{
	Use:   "rewind",
	Short: "Rewind current timer by number of minutes",
	Long:  "Rewinds the current timer by number of minutes and rewind the previous timer's stop time.",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		_, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("argument must be an integer: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		minutes, _ := strconv.Atoi(args[0])
		handleRewind(client, minutes)
	},
}

func init() {
	rootCmd.AddCommand(rewindCmd)
}

func handleRewind(client *toggl.TogglClient, minutes int) {
	entries, err := client.GetTimeEntries()
	if err != nil {
		log.Fatalln(err)
	}

	if len(entries) == 0 {
		log.Fatalln("no timers")
	}

	entries[0].Start = entries[0].Start.Add(time.Duration(-minutes) * time.Minute)

	// update this to use bulk edit PATCH endpoint instead
	err = client.UpdateTimeEntry(entries[0])

	if err != nil {
		log.Fatalln(err)
	}
	if len(entries) > 1 {
		newStop := entries[1].Stop.Add(time.Duration(-minutes) * time.Minute)
		entries[1].Stop = &newStop
		err := client.UpdateTimeEntry(entries[1])
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("updated start time")
}
