package api

import (
	"github.com/RootControl/twitch/entities"
)

const SEARCH_CHANNELS = "search/channels"

type ChannelResponse struct {
	Data []entities.Channel `json:"data"`
}

func GetChannel(broadcasterID string) (*ChannelResponse, error) {
	return getChannel(broadcasterID, nil)
}

func getChannel(broadcasterID string, e Executor) (*ChannelResponse, error) {
	resp, err := fetch[ChannelResponse](e, "channels", Q("broadcaster_id", broadcasterID))
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchChannelsResponse holds results from search/channels, which matches
// channel names and returns both live and offline broadcasters.
type SearchChannelsResponse struct {
	Data       []entities.SearchedChannel `json:"data"`
	Pagination entities.Pagination        `json:"pagination"`
}

// SearchChannels finds channels whose name matches query. When liveOnly is
// true only currently-live broadcasters are returned.
func SearchChannels(query string, liveOnly bool, limit int) (SearchChannelsResponse, error) {
	return searchChannels(query, liveOnly, limit, nil)
}

func searchChannels(query string, liveOnly bool, limit int, e Executor) (SearchChannelsResponse, error) {
	params := []Param{
		Q("query", query),
		Q("first", clampLimit(limit)),
	}
	if liveOnly {
		params = append(params, Q("live_only", "true"))
	}
	return fetch[SearchChannelsResponse](e, SEARCH_CHANNELS, params...)
}
