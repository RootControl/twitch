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
	RunE: func(cmd *cobra.Command, args []string) error {
		login := args[0]
		resp, err := api.GetUser(login)
		if err != nil {
			return err
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("user %q not found", login)
		}
		if err := api.FollowChannel(resp.Data[0].ID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Now following %s\n", resp.Data[0].DisplayName)
		return nil
	},
}

var unfollowCmd = &cobra.Command{
	Use:   "unfollow <username>",
	Short: "Unfollow a Twitch channel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		login := args[0]
		resp, err := api.GetUser(login)
		if err != nil {
			return err
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("user %q not found", login)
		}
		if err := api.UnfollowChannel(resp.Data[0].ID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Unfollowed %s\n", resp.Data[0].DisplayName)
		return nil
	},
}
