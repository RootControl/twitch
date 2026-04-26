package cmd

import (
	"fmt"
	"os"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var userJSON bool

var userCmd = &cobra.Command{
	Use:   "user <login>",
	Short: "Get info about a Twitch user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp := api.GetUser(args[0])
		if len(resp.Data) == 0 {
			fmt.Fprintf(os.Stderr, "User %q not found\n", args[0])
			os.Exit(1)
		}
		if userJSON {
			printJSON(resp.Data[0])
			return
		}
		fmt.Println(resp.Data[0].ToString())
	},
}

func init() {
	userCmd.Flags().BoolVar(&userJSON, "json", false, "Output result as JSON")
}
