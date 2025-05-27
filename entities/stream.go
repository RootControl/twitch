package entities

import "fmt"

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
	return &Stream{
		Id:           "",
		UserId:       "",
		UserLogin:    "",
		Username:     "",
		GameId:       "",
		GameName:     "",
		Type:         "",
		Title:        "",
		Tags:         make([]string, 0),
		ViewerCount:  0,
		StartedAt:    "",
		Language:     "",
		ThumbnailUrl: "",
		TagsIds:      make([]string, 0),
		IsMature:     false,
	}
}

func (s *Stream) GetMainInfo() string {
	return fmt.Sprintf("User: %s\nTitle: %s\nCategory: %s\n", s.Username, s.Title, s.GameName)
}
