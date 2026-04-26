package cmd

import (
	"fmt"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var followCmd = &cobra.Command{
	Use:   "follow <username>",
	Short: "Follow a Twitch channel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		login := args[0]
		resp := api.GetUser(login)
		if len(resp.Data) == 0 {
			fmt.Printf("User %q not found\n", login)
			return
		}
		api.FollowChannel(resp.Data[0].ID)
		fmt.Printf("Now following %s\n", resp.Data[0].DisplayName)
	},
}

var unfollowCmd = &cobra.Command{
	Use:   "unfollow <username>",
	Short: "Unfollow a Twitch channel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		login := args[0]
		resp := api.GetUser(login)
		if len(resp.Data) == 0 {
			fmt.Printf("User %q not found\n", login)
			return
		}
		api.UnfollowChannel(resp.Data[0].ID)
		fmt.Printf("Unfollowed %s\n", resp.Data[0].DisplayName)
	},
}
