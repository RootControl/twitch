package entities

import (
	"fmt"
	"time"
)

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
	TagsIds      []string `json:"tag_ids"`
	IsMature     bool     `json:"is_mature"`
}

func NewStream() *Stream {
	return &Stream{}
}

// URL is the twitch.tv address of the stream.
func (s *Stream) URL() string {
	return TwitchURL + s.UserLogin
}

// Uptime returns how long the stream has been live, formatted as "2h14m".
// It returns an empty string when started_at is missing or unparseable.
func (s *Stream) Uptime() string {
	return uptimeSince(s.StartedAt, time.Now())
}

func uptimeSince(startedAt string, now time.Time) string {
	if startedAt == "" {
		return ""
	}
	start, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		return ""
	}
	d := now.Sub(start)
	if d < 0 {
		return ""
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%02dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func (s *Stream) GetMainInfo() string {
	info := fmt.Sprintf("User: %s\nTitle: %s\nCategory: %s\n", s.Username, s.Title, s.GameName)
	if up := s.Uptime(); up != "" {
		info += fmt.Sprintf("Uptime: %s\n", up)
	}
	return info + fmt.Sprintf("Link: %s\n", s.URL())
}
