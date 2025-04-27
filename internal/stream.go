package internal

import (
	"encoding/json"
	"io"
	"net/http"
)

type Stream struct {
	Username string `json:"user_name"`
	Title    string `json:"title"`
}

type StreamResponse struct {
	Data []Stream `json:"data"`
}

const urlStreams = "https://api.twitch.tv/helix/streams"
const top = "?first=100"

func GetStreamers(accessToken string) ([]Stream, error) {
	auth := NewTwitchAuth()

	client := &http.Client{}

	req, err := http.NewRequest("GET", urlStreams+top, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-Id", auth.ClientId)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var streamResponse StreamResponse
	err = json.Unmarshal(bodyBytes, &streamResponse)
	if err != nil {
		return nil, err
	}

	return streamResponse.Data, nil
}
