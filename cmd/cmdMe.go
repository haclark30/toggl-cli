package cmd

import (
	"fmt"
	"log"

	"github.com/haclark30/toggl-cli/toggl"
	"github.com/spf13/cobra"
)

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Get info about current user.",
	Run: func(cmd *cobra.Command, args []string) {
		handleMe(client)
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}

func handleMe(client *toggl.TogglClient) {
	me, err := client.Me()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(me.Fullname)
}
