package entities

import "fmt"

type Channel struct {
	BroadcasterID       string `json:"broadcaster_id"`
	BroadcasterLogin    string `json:"broadcaster_login"`
	BroadcasterName     string `json:"broadcaster_name"`
	BroadcasterLanguage string `json:"broadcaster_language"`
	GameID              string `json:"game_id"`
	GameName            string `json:"game_name"`
	Title               string `json:"title"`
	Delay               int    `json:"delay"`
	Tags                []string `json:"tags"`
}

func (c *Channel) ToString() string {
	return fmt.Sprintf(
		"Channel:  %s\nTitle:    %s\nCategory: %s\nLanguage: %s\nLink:     %s",
		c.BroadcasterName, c.Title, c.GameName, c.BroadcasterLanguage,
		TwitchURL+c.BroadcasterLogin,
	)
}
