package cmd

import (
	"fmt"
	"os"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var channelJSON bool

var channelCmd = &cobra.Command{
	Use:   "channel <username>",
	Short: "Show info about a channel (title, category, language)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userResp := api.GetUser(args[0])
		if len(userResp.Data) == 0 {
			fmt.Fprintf(os.Stderr, "User %q not found\n", args[0])
			os.Exit(1)
		}
		broadcasterID := userResp.Data[0].ID
		chanResp := api.GetChannel(broadcasterID)
		if len(chanResp.Data) == 0 {
			fmt.Fprintf(os.Stderr, "Channel info not found for %q\n", args[0])
			os.Exit(1)
		}
		if channelJSON {
			printJSON(chanResp.Data[0])
			return
		}
		fmt.Println(chanResp.Data[0].ToString())
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelJSON, "json", false, "Output result as JSON")
}
