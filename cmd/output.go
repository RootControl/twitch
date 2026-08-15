package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/RootControl/twitch/entities"
	"golang.org/x/term"
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

// fallbackWidth is used when stdout is not a terminal or its size is unknown.
const fallbackWidth = 80

// terminalWidth reports the width of the terminal attached to stdout, in
// columns.
func terminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 20 {
		return fallbackWidth
	}
	return w
}

// runeWidth reports how many terminal columns r occupies. It is an
// approximation: it covers the wide and zero-width ranges that show up in
// Twitch titles (emoji, CJK, combining marks) and treats everything else as a
// single column. Over-estimating is safe here — it only truncates earlier.
func runeWidth(r rune) int {
	switch {
	case r == 0x200D, r >= 0xFE00 && r <= 0xFE0F:
		return 0 // zero-width joiner, variation selectors
	case unicode.Is(unicode.Mn, r), unicode.Is(unicode.Me, r):
		return 0 // combining marks
	case r >= 0x1100 && r <= 0x115F, // Hangul Jamo
		r >= 0x2E80 && r <= 0xA4CF, // CJK radicals through Yi
		r >= 0xAC00 && r <= 0xD7A3, // Hangul syllables
		r >= 0xF900 && r <= 0xFAFF, // CJK compatibility ideographs
		r >= 0xFE30 && r <= 0xFE6F, // CJK compatibility forms
		r >= 0xFF00 && r <= 0xFF60, // fullwidth forms
		r >= 0xFFE0 && r <= 0xFFE6,
		r >= 0x1F300 && r <= 0x1F64F, // emoji, pictographs, emoticons
		r >= 0x1F680 && r <= 0x1F6FF, // transport and map symbols
		r >= 0x1F900 && r <= 0x1F9FF, // supplemental symbols
		r >= 0x1FA70 && r <= 0x1FAFF, // extended-A symbols
		r >= 0x20000 && r <= 0x3FFFD: // CJK extension planes
		return 2
	default:
		return 1
	}
}

// displayWidth reports how many terminal columns s occupies.
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

// truncate shortens s to at most max terminal columns, marking a cut with an
// ellipsis. It measures display width rather than rune count, so a title full
// of emoji does not overflow the line it was supposed to fit on.
func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if displayWidth(s) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}

	limit := max - 1 // leave a column for the ellipsis
	var b strings.Builder
	w := 0
	for _, r := range s {
		rw := runeWidth(r)
		if w+rw > limit {
			break
		}
		b.WriteRune(r)
		w += rw
	}
	return b.String() + "…"
}
