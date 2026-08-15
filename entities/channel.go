package entities

import "fmt"

type Channel struct {
	BroadcasterID       string   `json:"broadcaster_id"`
	BroadcasterLogin    string   `json:"broadcaster_login"`
	BroadcasterName     string   `json:"broadcaster_name"`
	BroadcasterLanguage string   `json:"broadcaster_language"`
	GameID              string   `json:"game_id"`
	GameName            string   `json:"game_name"`
	Title               string   `json:"title"`
	Delay               int      `json:"delay"`
	Tags                []string `json:"tags"`
}

func (c *Channel) ToString() string {
	return fmt.Sprintf(
		"Channel:  %s\nTitle:    %s\nCategory: %s\nLanguage: %s\nLink:     %s",
		c.BroadcasterName, c.Title, c.GameName, c.BroadcasterLanguage,
		TwitchURL+c.BroadcasterLogin,
	)
}

// SearchedChannel is a result from the search/channels endpoint. It differs
// from Channel: it carries live state but not the channel's delay or tags in
// the same shape.
type SearchedChannel struct {
	BroadcasterID       string   `json:"id"`
	BroadcasterLogin    string   `json:"broadcaster_login"`
	BroadcasterName     string   `json:"display_name"`
	BroadcasterLanguage string   `json:"broadcaster_language"`
	GameID              string   `json:"game_id"`
	GameName            string   `json:"game_name"`
	Title               string   `json:"title"`
	IsLive              bool     `json:"is_live"`
	StartedAt           string   `json:"started_at"`
	ThumbnailURL        string   `json:"thumbnail_url"`
	Tags                []string `json:"tags"`
}

func (c *SearchedChannel) ToString() string {
	status := "offline"
	if c.IsLive {
		status = "LIVE"
	}
	return fmt.Sprintf(
		"%-8s %s\n         %s\n         %s\n         %s",
		status, c.BroadcasterName, c.Title, c.GameName, TwitchURL+c.BroadcasterLogin,
	)
}
