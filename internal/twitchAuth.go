package internal

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type TwitchAuth struct {
	ClientId     string
	ClientSecret string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

const urlToken = "https://id.twitch.tv/oauth2/token"

func NewTwitchAuth() TwitchAuth {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return TwitchAuth{
		ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
	}
}

func GetAccessToken() (string, error) {
	auth := NewTwitchAuth()

	form := url.Values{}
	form.Add("client_id", auth.ClientId)
	form.Add("client_secret", auth.ClientSecret)
	form.Add("grant_type", "client_credentials")

	resp, err := http.Post(
		urlToken,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.AccessToken, nil
}
