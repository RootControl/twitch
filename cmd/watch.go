package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Poll followed streams and notify when someone goes live",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Watching followed streams (polling every %ds). Press Ctrl+C to stop.\n", watchInterval)
		live := map[string]bool{}

		for {
			resp := api.GetFollowedStreams()
			current := map[string]bool{}

			for _, s := range resp.Data {
				current[s.UserLogin] = true
				if !live[s.UserLogin] {
					msg := fmt.Sprintf("%s just went live: %s (%s)", s.Username, s.Title, s.GameName)
					notify(s.Username+" is live!", msg)
					fmt.Println(msg)
				}
			}

			live = current
			time.Sleep(time.Duration(watchInterval) * time.Second)
		}
	},
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 60, "Polling interval in seconds")
}

func notify(title, message string) {
	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, message, title)
		exec.Command("osascript", "-e", script).Run()
	case "linux":
		exec.Command("notify-send", title, message).Run()
	case "windows":
		// Windows toast notifications require PowerShell or a third-party lib; skip silently.
	}
}
