package cmd

import (
	"fmt"

	"github.com/RootControl/twitch/api"
	"github.com/RootControl/twitch/entities"
	"github.com/spf13/cobra"
)

var (
	searchJSON     bool
	searchLiveOnly bool
	searchLimit    int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search channels by name (live and offline)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.SearchChannels(args[0], searchLiveOnly, searchLimit)
		if err != nil {
			return err
		}

		if searchJSON {
			if resp.Data == nil {
				resp.Data = []entities.SearchedChannel{}
			}
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("no channels found for %q", args[0])
		}
		for i := range resp.Data {
			fmt.Fprintln(cmd.OutOrStdout(), resp.Data[i].ToString())
			fmt.Fprintln(cmd.OutOrStdout())
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output results as JSON")
	searchCmd.Flags().BoolVarP(&searchLiveOnly, "live", "l", false, "Only show channels that are currently live")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 20, "Maximum number of channels to fetch (1-100)")
}
