package cmd

import (
	"fmt"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var channelJSON bool

var channelCmd = &cobra.Command{
	Use:   "channel <username>",
	Short: "Show info about a channel (title, category, language)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userResp, err := api.GetUser(args[0])
		if err != nil {
			return err
		}
		if len(userResp.Data) == 0 {
			return fmt.Errorf("user %q not found", args[0])
		}

		chanResp, err := api.GetChannel(userResp.Data[0].ID)
		if err != nil {
			return err
		}
		if len(chanResp.Data) == 0 {
			return fmt.Errorf("channel info not found for %q", args[0])
		}

		if channelJSON {
			return printJSON(cmd.OutOrStdout(), chanResp.Data[0])
		}
		fmt.Fprintln(cmd.OutOrStdout(), chanResp.Data[0].ToString())
		return nil
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelJSON, "json", false, "Output result as JSON")
}
