package cmd

import (
	"fmt"
	"os"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var categoriesJSON bool

var categoriesCmd = &cobra.Command{
	Use:   "categories <query>",
	Short: "Search for Twitch categories/games",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp := api.GetCategories(args[0])
		if len(resp.Data) == 0 {
			fmt.Fprintf(os.Stderr, "No categories found for %q\n", args[0])
			os.Exit(1)
		}
		if categoriesJSON {
			printJSON(resp.Data)
			return
		}
		for _, c := range resp.Data {
			fmt.Println(c.ToString())
		}
	},
}

func init() {
	categoriesCmd.Flags().BoolVar(&categoriesJSON, "json", false, "Output results as JSON")
}
