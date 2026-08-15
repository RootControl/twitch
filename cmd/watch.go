package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/RootControl/twitch/api"
	"github.com/spf13/cobra"
)

var (
	watchInterval int
	watchNotify   bool
)

// minWatchInterval keeps the poller from hammering the API (and from spinning
// in a tight loop when --interval 0 is passed).
const minWatchInterval = 10

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Poll followed streams and notify when someone goes live",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if watchInterval < minWatchInterval {
			return fmt.Errorf("--interval must be at least %d seconds, got %d", minWatchInterval, watchInterval)
		}

		out := cmd.OutOrStdout()

		// Seed the baseline with whoever is already live, so the first poll
		// does not announce every ongoing stream as if it just started.
		live := map[string]bool{}
		resp, err := api.GetFollowedStreams(api.DefaultLimit)
		if err != nil {
			return err
		}
		for _, s := range resp.Data {
			live[s.UserLogin] = true
		}
		fmt.Fprintf(out, "Watching followed streams (%d live now, polling every %ds). Press Ctrl+C to stop.\n",
			len(live), watchInterval)

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		ticker := time.NewTicker(time.Duration(watchInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Fprintln(out, "Stopped watching.")
				return nil
			case <-ticker.C:
			}

			resp, err := api.GetFollowedStreams(api.DefaultLimit)
			if err != nil {
				// A transient failure must not kill a long-running watch.
				fmt.Fprintf(cmd.ErrOrStderr(), "poll failed, retrying in %ds: %v\n", watchInterval, err)
				continue
			}

			current := make(map[string]bool, len(resp.Data))
			for _, s := range resp.Data {
				current[s.UserLogin] = true
				if !live[s.UserLogin] {
					msg := fmt.Sprintf("%s just went live: %s (%s)", s.Username, s.Title, s.GameName)
					if watchNotify {
						notify(s.Username+" is live!", msg)
					}
					fmt.Fprintf(out, "[%s] %s\n", time.Now().Format("15:04:05"), msg)
				}
			}
			live = current
		}
	},
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 60, "Polling interval in seconds (minimum 10)")
	watchCmd.Flags().BoolVar(&watchNotify, "notify", true, "Send a desktop notification when a channel goes live")
}

func notify(title, message string) {
	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf("display notification %s with title %s",
			appleScriptString(message), appleScriptString(title))
		exec.Command("osascript", "-e", script).Run()
	case "linux":
		exec.Command("notify-send", title, message).Run()
	case "windows":
		// Windows toast notifications require PowerShell or a third-party lib; skip silently.
	}
}

// appleScriptString quotes a Go string as an AppleScript literal. Stream
// titles routinely contain double quotes and backslashes, which would
// otherwise terminate the literal and break (or alter) the script.
func appleScriptString(s string) string {
	quoted := make([]rune, 0, len(s)+2)
	quoted = append(quoted, '"')
	for _, r := range s {
		if r == '"' || r == '\\' {
			quoted = append(quoted, '\\')
		}
		quoted = append(quoted, r)
	}
	quoted = append(quoted, '"')
	return string(quoted)
}
