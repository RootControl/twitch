package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/RootControl/twitch/entities"
)

func printJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("json encode error: %w", err)
	}
	return nil
}

// isInteractive reports whether stdin and stdout are both attached to a
// terminal. The interactive picker cannot run otherwise (in a pipe it fails
// with an opaque error), so commands fall back to plain listings.
func isInteractive() bool {
	return isCharDevice(os.Stdin) && isCharDevice(os.Stdout)
}

func isCharDevice(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// printStreamTable writes streams as an aligned table.
func printStreamTable(w io.Writer, streams []entities.Stream) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "CHANNEL\tVIEWERS\tUPTIME\tCATEGORY\tTITLE")
	for i := range streams {
		s := &streams[i]
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\n",
			s.Username, s.ViewerCount, dashIfEmpty(s.Uptime()), dashIfEmpty(s.GameName), truncate(s.Title, 60))
	}
	tw.Flush()
}

func dashIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 1 {
		return string(runes[:max])
	}
	return string(runes[:max-1]) + "…"
}
