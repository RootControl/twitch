package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RootControl/twitch/entities"
)

const (
	GET_LIVE_STREAMS     = "streams"
	GET_FOLLOWED_STREAMS = "streams/followed"
)

// DefaultLimit is the page size used when a caller does not specify one.
const DefaultLimit = 30

// maxLimit is the largest page size the Helix API accepts.
const maxLimit = 100

type StreamResponse struct {
	Data       []entities.Stream   `json:"data"`
	Pagination entities.Pagination `json:"pagination"`
}

func NewStreamResponse() *StreamResponse {
	return &StreamResponse{
		Data:       make([]entities.Stream, 0),
		Pagination: entities.Pagination{},
	}
}

// clampLimit keeps the requested page size within what Helix accepts.
func clampLimit(limit int) string {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return strconv.Itoa(limit)
}

func GetLiveStreams(limit int) (StreamResponse, error) {
	return getLiveStreams(limit, nil)
}

func getLiveStreams(limit int, e Executor) (StreamResponse, error) {
	return fetch[StreamResponse](e, GET_LIVE_STREAMS,
		Q("type", "live"),
		Q("first", clampLimit(limit)),
	)
}

func GetFollowedStreams(limit int) (StreamResponse, error) {
	return getFollowedStreams(limit, nil)
}

func getFollowedStreams(limit int, e Executor) (StreamResponse, error) {
	userID, err := resolveUserID(e)
	if err != nil {
		return StreamResponse{}, err
	}
	return fetch[StreamResponse](e, GET_FOLLOWED_STREAMS,
		Q("user_id", userID),
		Q("first", clampLimit(limit)),
	)
}

func GetStreamsByGame(gameName string, limit int) (StreamResponse, error) {
	return getStreamsByGame(gameName, limit, nil)
}

func getStreamsByGame(gameName string, limit int, e Executor) (StreamResponse, error) {
	// Resolve the game name to an ID first.
	cats, err := getCategories(gameName, 0, e)
	if err != nil {
		return StreamResponse{}, err
	}
	if len(cats.Data) == 0 {
		return StreamResponse{}, fmt.Errorf("no category matches %q", gameName)
	}

	gameID := bestCategoryMatch(cats.Data, gameName).ID

	return fetch[StreamResponse](e, GET_LIVE_STREAMS,
		Q("type", "live"),
		Q("game_id", gameID),
		Q("first", clampLimit(limit)),
	)
}

// bestCategoryMatch prefers an exact (case-insensitive) name match over the
// first fuzzy result, so `--game Rust` does not resolve to "Rust Console
// Edition" just because the search ranked it first.
func bestCategoryMatch(cats []entities.Category, want string) entities.Category {
	for _, c := range cats {
		if strings.EqualFold(c.Name, want) {
			return c
		}
	}
	return cats[0]
}
