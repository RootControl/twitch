package cmd

import (
	"fmt"
	"os"

	"github.com/RootControl/twitch/api"
	"github.com/RootControl/twitch/browser"
	"github.com/RootControl/twitch/entities"
	"github.com/RootControl/twitch/player"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var streamsGame string
var streamsJSON bool

var streamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "List and open top live streams",
	Run: func(cmd *cobra.Command, args []string) {
		var resp api.StreamResponse
		if streamsGame != "" {
			resp = api.GetStreamsByGame(streamsGame)
			if len(resp.Data) == 0 {
				fmt.Fprintf(os.Stderr, "No live streams found for game %q\n", streamsGame)
				os.Exit(1)
			}
		} else {
			resp = api.GetLiveStreams()
		}
		if streamsJSON {
			printJSON(resp.Data)
			return
		}
		pickAndOpen(resp.Data)
	},
}

var followedJSON bool

var followedCmd = &cobra.Command{
	Use:   "followed",
	Short: "List and open followed live streams",
	Run: func(cmd *cobra.Command, args []string) {
		resp := api.GetFollowedStreams()
		if len(resp.Data) == 0 {
			fmt.Println("No followed streams are live right now.")
			return
		}
		if followedJSON {
			printJSON(resp.Data)
			return
		}
		pickAndOpen(resp.Data)
	},
}

func init() {
	streamsCmd.Flags().StringVarP(&streamsGame, "game", "g", "", "Filter streams by game/category name")
	streamsCmd.Flags().BoolVar(&streamsJSON, "json", false, "Output results as JSON")
	followedCmd.Flags().BoolVar(&followedJSON, "json", false, "Output results as JSON")
}

func pickAndOpen(streams []entities.Stream) {
	if len(streams) == 0 {
		fmt.Println("No streams found.")
		return
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
{{ "Viewers:" | faint }}  {{ .ViewerCount }}`,
	}

	prompt := promptui.Select{
		Label:     "Select a stream",
		Items:     streams,
		Templates: templates,
		Size:      15,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return
	}

	stream := streams[i]
	url := entities.TwitchURL + stream.UserLogin
	openStream(stream.Username, url)
}

func openStream(name, url string) {
	players := player.Available()
	if len(players) == 0 {
		fmt.Printf("Opening %s in browser\n", url)
		if err := browser.OpenBrowser(url); err != nil {
			fmt.Fprintf(os.Stderr, "Could not open browser: %v\n", err)
		}
		return
	}

	// Build choice list: available players + browser fallback.
	type choice struct{ label, kind string; p *player.Player }
	choices := make([]choice, 0, len(players)+1)
	for idx := range players {
		p := players[idx]
		choices = append(choices, choice{label: p.Name, kind: "player", p: &p})
	}
	choices = append(choices, choice{label: "browser", kind: "browser"})

	labels := make([]string, len(choices))
	for i, c := range choices {
		labels[i] = c.label
	}

	sel := promptui.Select{
		Label: fmt.Sprintf("Open %s with", name),
		Items: labels,
	}
	idx, _, err := sel.Run()
	if err != nil {
		return
	}

	picked := choices[idx]
	if picked.kind == "browser" {
		fmt.Printf("Opening %s in browser\n", url)
		if err := browser.OpenBrowser(url); err != nil {
			fmt.Fprintf(os.Stderr, "Could not open browser: %v\n", err)
		}
	} else {
		fmt.Printf("Opening %s with %s\n", name, picked.p.Name)
		if err := picked.p.Open(url); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}
