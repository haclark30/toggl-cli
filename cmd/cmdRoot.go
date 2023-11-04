package cmd

import (
	"fmt"
	"os"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/spf13/cobra"
)

var client *toggl.TogglClient

const workspaceId = 7636849

var rootCmd = &cobra.Command{
	Use:   "tempus",
	Short: "Command line toggl client",
	Long:  "A command-line time tracking app using Toggl as a backend.",
}

func init() {
	apiToken := getApiToken()
	client = toggl.NewTogglClient(apiToken)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
