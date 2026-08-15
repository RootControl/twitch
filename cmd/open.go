package cmd

import (
	"fmt"

	"github.com/RootControl/twitch/api"
	"github.com/RootControl/twitch/entities"
	"github.com/spf13/cobra"
)

var openPlayer string

var openCmd = &cobra.Command{
	Use:   "open <username>",
	Short: "Open a channel in a player or the browser",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		login := args[0]

		// Verify the channel exists before launching anything, so a typo
		// produces a clear message instead of a player error or a 404 page.
		resp, err := api.GetUser(login)
		if err != nil {
			return err
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("user %q not found", login)
		}

		user := resp.Data[0]
		name := user.DisplayName
		if name == "" {
			name = user.Login
		}
		return openStream(cmd, name, entities.TwitchURL+user.Login, openPlayer)
	},
}

func init() {
	openCmd.Flags().StringVarP(&openPlayer, "player", "p", "", "Open with this player (streamlink, mpv, browser)")
}
