package cmd

import (
	"fmt"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var (
	categoriesJSON  bool
	categoriesLimit int
	topJSON         bool
	topLimit        int
)

var categoriesCmd = &cobra.Command{
	Use:   "categories <query>",
	Short: "Search for Twitch categories/games",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.GetCategories(args[0], categoriesLimit)
		if err != nil {
			return err
		}
		if categoriesJSON {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		if len(resp.Data) == 0 {
			return fmt.Errorf("no categories found for %q", args[0])
		}
		for i := range resp.Data {
			fmt.Fprintln(cmd.OutOrStdout(), resp.Data[i].ToString())
		}
		return nil
	},
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "List the most-watched categories right now",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.GetTopGames(topLimit)
		if err != nil {
			return err
		}
		if topJSON {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		if len(resp.Data) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No categories returned.")
			return nil
		}
		for i, c := range resp.Data {
			fmt.Fprintf(cmd.OutOrStdout(), "%2d. %s\n", i+1, c.Name)
		}
		return nil
	},
}

func init() {
	categoriesCmd.Flags().BoolVar(&categoriesJSON, "json", false, "Output results as JSON")
	categoriesCmd.Flags().IntVarP(&categoriesLimit, "limit", "n", 50, "Maximum number of categories to fetch (1-100)")

	topCmd.Flags().BoolVar(&topJSON, "json", false, "Output results as JSON")
	topCmd.Flags().IntVarP(&topLimit, "limit", "n", 20, "Maximum number of categories to fetch (1-100)")
}
