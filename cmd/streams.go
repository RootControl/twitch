package cmd

import (
	"fmt"
	"strings"

	"github.com/RootControl/twitch/api"
	"github.com/RootControl/twitch/browser"
	"github.com/RootControl/twitch/entities"
	"github.com/RootControl/twitch/player"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	streamsGame   string
	streamsJSON   bool
	streamsLimit  int
	streamsPlayer string

	followedJSON   bool
	followedLimit  int
	followedPlayer string
)

var streamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "List and open top live streams",
	RunE: func(cmd *cobra.Command, args []string) error {
		var resp api.StreamResponse
		var err error

		if streamsGame != "" {
			resp, err = api.GetStreamsByGame(streamsGame, streamsLimit)
		} else {
			resp, err = api.GetLiveStreams(streamsLimit)
		}
		if err != nil {
			return err
		}

		if streamsJSON {
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		if len(resp.Data) == 0 {
			if streamsGame != "" {
				return fmt.Errorf("no live streams found for game %q", streamsGame)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "No live streams found.")
			return nil
		}
		return pickAndOpen(cmd, resp.Data, streamsPlayer)
	},
}

var followedCmd = &cobra.Command{
	Use:   "followed",
	Short: "List and open followed live streams",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.GetFollowedStreams(followedLimit)
		if err != nil {
			return err
		}

		if followedJSON {
			// Always emit valid JSON, including for an empty result set.
			if resp.Data == nil {
				resp.Data = []entities.Stream{}
			}
			return printJSON(cmd.OutOrStdout(), resp.Data)
		}
		if len(resp.Data) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No followed streams are live right now.")
			return nil
		}
		return pickAndOpen(cmd, resp.Data, followedPlayer)
	},
}

func init() {
	streamsCmd.Flags().StringVarP(&streamsGame, "game", "g", "", "Filter streams by game/category name")
	streamsCmd.Flags().BoolVar(&streamsJSON, "json", false, "Output results as JSON")
	streamsCmd.Flags().IntVarP(&streamsLimit, "limit", "n", api.DefaultLimit, "Maximum number of streams to fetch (1-100)")
	streamsCmd.Flags().StringVarP(&streamsPlayer, "player", "p", "", "Open directly with this player (streamlink, mpv, browser)")

	followedCmd.Flags().BoolVar(&followedJSON, "json", false, "Output results as JSON")
	followedCmd.Flags().IntVarP(&followedLimit, "limit", "n", api.DefaultLimit, "Maximum number of streams to fetch (1-100)")
	followedCmd.Flags().StringVarP(&followedPlayer, "player", "p", "", "Open directly with this player (streamlink, mpv, browser)")
}

func pickAndOpen(cmd *cobra.Command, streams []entities.Stream, preferred string) error {
	out := cmd.OutOrStdout()
	if len(streams) == 0 {
		fmt.Fprintln(out, "No streams found.")
		return nil
	}

	// Without a terminal the picker cannot run, so just print the results.
	if !isInteractive() {
		printStreamTable(out, streams)
		return nil
	}

	items := newStreamItems(streams, terminalWidth())
	templates := streamSelectTemplates()

	prompt := promptui.Select{
		Label:     "Select a stream",
		Items:     items,
		Templates: templates,
		Size:      15,
		// Search the underlying stream so matching uses the full text, not
		// the truncated display copy.
		Searcher: func(input string, index int) bool {
			return matchesStream(&items[index].stream, input)
		},
		StartInSearchMode: false,
	}

	i, _, err := prompt.Run()
	if err != nil {
		// A cancelled prompt (Ctrl+C / Esc) is a normal exit, not an error.
		return nil
	}

	stream := items[i].stream
	return openStream(cmd, stream.Username, stream.URL(), preferred)
}

// streamSelectTemplates renders a streamItem. Every field it interpolates is
// pre-truncated by newStreamItems, so no rendered line can exceed the terminal
// width — see the comment on streamItem for why that matters.
func streamSelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Name | cyan }} ({{ .Game | yellow }}) - {{ .Viewers }} viewers",
		Inactive: "  {{ .Name }} ({{ .Game }}) - {{ .Viewers }} viewers",
		Selected: "▶ {{ .Name | cyan }}",
		// The padding belongs inside the quoted labels: whitespace within a
		// template action is syntax and is not emitted, so spacing it there
		// leaves the columns ragged.
		Details: `
--- Stream ---
{{ "Channel:" | faint }}  {{ .Name }}
{{ "Title:  " | faint }}  {{ .Title }}
{{ "Game:   " | faint }}  {{ .Game }}
{{ "Viewers:" | faint }}  {{ .Viewers }}
{{ "Uptime: " | faint }}  {{ .Uptime }}`,
	}
}

// streamItem is a display-ready view of a stream, with every field already cut
// to fit the terminal.
//
// This exists because promptui's screen buffer counts the *logical* lines it
// writes, then moves the cursor up that many rows to redraw. A line wider than
// the terminal wraps onto extra rows that the buffer does not know about, so
// the redraw erases too little and leaves the previous frame on screen. Long
// stream titles are more than enough to trigger it.
type streamItem struct {
	Name    string
	Game    string
	Title   string
	Uptime  string
	Viewers int

	stream entities.Stream
}

func newStreamItems(streams []entities.Stream, width int) []streamItem {
	nameWidth, gameWidth, detailWidth := itemWidths(width)

	items := make([]streamItem, len(streams))
	for i, s := range streams {
		items[i] = streamItem{
			Name:    truncate(s.Username, nameWidth),
			Game:    truncate(s.GameName, gameWidth),
			Title:   truncate(s.Title, detailWidth),
			Uptime:  s.Uptime(),
			Viewers: s.ViewerCount,
			stream:  s,
		}
	}
	return items
}

// itemWidths splits the terminal width into per-field budgets.
func itemWidths(width int) (name, game, detail int) {
	// A list row renders as "▶ NAME (GAME) - 12345 viewers"; everything
	// outside the two names costs about 23 columns.
	const listOverhead = 23
	// A detail row renders as "Channel:  VALUE" — a 10-column label, plus a
	// column of slack so the cursor never lands in the final cell.
	const labelWidth = 11

	avail := max(width-listOverhead, 12)
	name = avail / 2
	game = avail - name
	if name > 24 {
		game += name - 24
		name = 24
	}
	game = min(game, 32)

	detail = max(width-labelWidth, 20)
	return name, game, detail
}

func matchesStream(s *entities.Stream, input string) bool {
	input = strings.ToLower(input)
	return strings.Contains(strings.ToLower(s.Username), input) ||
		strings.Contains(strings.ToLower(s.GameName), input) ||
		strings.Contains(strings.ToLower(s.Title), input)
}

// openStream launches the URL, either with the preferred player or by asking.
func openStream(cmd *cobra.Command, name, url, preferred string) error {
	out := cmd.OutOrStdout()
	players := player.Available()

	if preferred != "" {
		if isBrowserChoice(preferred) {
			return openInBrowser(cmd, url)
		}
		for _, p := range players {
			if strings.EqualFold(p.Name, preferred) {
				fmt.Fprintf(out, "Opening %s with %s\n", name, p.Name)
				return p.Open(url)
			}
		}
		return fmt.Errorf("player %q is not available (installed: %s)", preferred, playerNames(players))
	}

	if len(players) == 0 || !isInteractive() {
		return openInBrowser(cmd, url)
	}

	labels := make([]string, 0, len(players)+1)
	for _, p := range players {
		labels = append(labels, p.Name)
	}
	labels = append(labels, "browser")

	sel := promptui.Select{
		Label: fmt.Sprintf("Open %s with", name),
		Items: labels,
	}
	idx, _, err := sel.Run()
	if err != nil {
		return nil // cancelled
	}

	if idx == len(players) {
		return openInBrowser(cmd, url)
	}
	p := players[idx]
	fmt.Fprintf(out, "Opening %s with %s\n", name, p.Name)
	return p.Open(url)
}

func openInBrowser(cmd *cobra.Command, url string) error {
	fmt.Fprintf(cmd.OutOrStdout(), "Opening %s in browser\n", url)
	if err := browser.OpenBrowser(url); err != nil {
		return fmt.Errorf("could not open browser: %w", err)
	}
	return nil
}

func isBrowserChoice(s string) bool {
	return strings.EqualFold(s, "browser")
}

func playerNames(players []player.Player) string {
	if len(players) == 0 {
		return "none"
	}
	names := make([]string, 0, len(players))
	for _, p := range players {
		names = append(names, p.Name)
	}
	names = append(names, "browser")
	return strings.Join(names, ", ")
}
