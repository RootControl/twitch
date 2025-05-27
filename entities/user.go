package entities

import "fmt"

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) ToString() string {
	return fmt.Sprintf("ID: %s\tLogin: %s\nDisplay Name: %s", u.ID, u.Login, u.DisplayName)
}
