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

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Username | cyan }} ({{ .GameName | yellow }}) - {{ .ViewerCount }} viewers",
		Inactive: "  {{ .Username }} ({{ .GameName }}) - {{ .ViewerCount }} viewers",
		Selected: "▶ {{ .Username | cyan }}",
		Details: `
--- Stream ---
{{ "Channel:" | faint }}  {{ .Username }}
{{ "Title:"   | faint }}  {{ .Title }}
{{ "Game:"    | faint }}  {{ .GameName }}
{{ "Viewers:" | faint }}  {{ .ViewerCount }}
{{ "Uptime:"  | faint }}  {{ .Uptime }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a stream",
		Items:     streams,
		Templates: templates,
		Size:      15,
		Searcher: func(input string, index int) bool {
			return matchesStream(&streams[index], input)
		},
		StartInSearchMode: false,
	}

	i, _, err := prompt.Run()
	if err != nil {
		// A cancelled prompt (Ctrl+C / Esc) is a normal exit, not an error.
		return nil
	}

	stream := streams[i]
	return openStream(cmd, stream.Username, stream.URL(), preferred)
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
