package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// defaultName is used when the invoked name cannot be determined. It is also
// the name rootCmd carries in tests, which drive rootCmd.Execute() directly
// rather than going through Execute().
const defaultName = "twitch"

var rootCmd = &cobra.Command{
	Use:           defaultName,
	Short:         "Twitch CLI wrapper for browsing streams, users, and categories",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	// Adopt the name the binary was actually invoked as, so usage text and
	// generated completion scripts match whatever it was installed as. This
	// matters because the default build name collides with the Twitch CLI
	// this tool shells out to, so it is commonly installed under another.
	rootCmd.Use = executableName()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// executableName reports the base name the process was invoked under. It uses
// os.Args[0] rather than os.Executable() so that invoking through a symlink
// reports the name the user typed, not the link's target.
func executableName() string {
	if len(os.Args) == 0 {
		return defaultName
	}

	name := filepath.Base(os.Args[0])
	name = strings.TrimSuffix(name, ".exe")

	// filepath.Base never returns "", but it does return "." or a separator
	// for degenerate input, and neither is usable as a command name.
	if name == "." || name == string(filepath.Separator) {
		return defaultName
	}
	return name
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
