package api

import (
	"encoding/json"
	"log"

	"github.com/RootControl/twitch/entities"
)

const (
	GET_LIVE_STREAMS = "streams"
)

type StreamResponse struct {
	Data       []entities.Stream `json:"data"`
	Pagination entities.Pagination
}

func NewStreamResponse() *StreamResponse {
	return &StreamResponse{
		Data:       make([]entities.Stream, 0),
		Pagination: entities.Pagination{},
	}
}

func GetLiveStreams() StreamResponse {
	request := NewApiRequest()

	responseBuffer := request.Get(GET_LIVE_STREAMS, "-q type=live", "-q first=30")

	var response StreamResponse
	err := json.Unmarshal(responseBuffer.Bytes(), &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return response
}
