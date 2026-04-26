package api

import (
	"encoding/json"
	"log"

	"github.com/RootControl/twitch/config"
	"github.com/RootControl/twitch/entities"
)

const (
	GET_LIVE_STREAMS     = "streams"
	GET_FOLLOWED_STREAMS = "streams/followed"
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
	return getLiveStreams(nil)
}

func getLiveStreams(e Executor) StreamResponse {
	request := newRequestWithExecutor(orShell(e))
	buf := request.Get(GET_LIVE_STREAMS, "-q type=live", "-q first=30")
	var response StreamResponse
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return response
}

func GetFollowedStreams() StreamResponse {
	return getFollowedStreams(nil)
}

func getFollowedStreams(e Executor) StreamResponse {
	request := newRequestWithExecutor(orShell(e))
	buf := request.Get(GET_FOLLOWED_STREAMS, "-q user_id="+config.MustUserID(), "-q first=50")
	var response StreamResponse
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return response
}

func GetStreamsByGame(gameName string) StreamResponse {
	return getStreamsByGame(gameName, nil)
}

func getStreamsByGame(gameName string, e Executor) StreamResponse {
	// Resolve game name to ID first.
	cats := getCategories(gameName, e)
	if len(cats.Data) == 0 {
		return StreamResponse{}
	}
	gameID := cats.Data[0].ID

	request := newRequestWithExecutor(orShell(e))
	buf := request.Get(GET_LIVE_STREAMS, "-q type=live", "-q game_id="+gameID, "-q first=30")
	var response StreamResponse
	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	return response
}

func orShell(e Executor) Executor {
	if e != nil {
		return e
	}
	return shellExecutor
}
