package entities

import "fmt"

const TwitchURL = "https://www.twitch.tv/"

type Stream struct {
	Id           string   `json:"id"`
	UserId       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	Username     string   `json:"user_name"`
	GameId       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Tags         []string `json:"tags"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	TagsIds      []string `json:"tags_ids"`
	IsMature     bool     `json:"is_mature"`
}

func NewStream() *Stream {
	return &Stream{}
}

func (s *Stream) GetMainInfo() string {
	return fmt.Sprintf("User: %s\nTitle: %s\nCategory: %s\nLink: %s\n", s.Username, s.Title, s.GameName, TwitchURL+s.UserLogin)
}
