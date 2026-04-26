package api

import (
	"encoding/json"
	"log"

	"github.com/RootControl/twitch/entities"
)

type ChannelResponse struct {
	Data []entities.Channel `json:"data"`
}

func GetChannel(broadcasterID string) *ChannelResponse {
	return getChannel(broadcasterID, nil)
}

func getChannel(broadcasterID string, e Executor) *ChannelResponse {
	request := newRequestWithExecutor(orShell(e))
	buf := request.Get("channels", "-q broadcaster_id="+broadcasterID)
	var response ChannelResponse
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return &response
}
