package api

import (
	"encoding/json"
	"log"

	"github.com/RootControl/twitch/entities"
)

const (
	GET_LIVE_STREAMS = "streams"
)

func GetLiveStreams() entities.StreamResponse {
	request := NewApiRequest()

	responseBuffer := request.Get(GET_LIVE_STREAMS, "-q type=live", "-q first=30")

	var response entities.StreamResponse
	err := json.Unmarshal(responseBuffer.Bytes(), &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return response
}
