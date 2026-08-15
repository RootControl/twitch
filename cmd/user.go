package cmd

import (
	"fmt"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var userJSON bool

var userCmd = &cobra.Command{
	Use:   "user <login>",
	Short: "Get info about a Twitch user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.GetUser(args[0])
		if err != nil {
			return err
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("user %q not found", args[0])
		}
		if userJSON {
			return printJSON(cmd.OutOrStdout(), resp.Data[0])
		}
		fmt.Fprintln(cmd.OutOrStdout(), resp.Data[0].ToString())
		return nil
	},
}

func init() {
	userCmd.Flags().BoolVar(&userJSON, "json", false, "Output result as JSON")
}
