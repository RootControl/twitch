package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "twitch",
	Short:         "Twitch CLI wrapper for browsing streams, users, and categories",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(streamsCmd)
	rootCmd.AddCommand(followedCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(channelCmd)
	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(topCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(followCmd)
	rootCmd.AddCommand(unfollowCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(completionCmd)
}
