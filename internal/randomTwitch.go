package internal

import (
	"fmt"
	"math/rand"
	"time"
)

type RandomTwitch struct {
	Streamers []Stream
}

const twichUrl = "https://www.twitch.tv/"

func NewRandomTwitch() RandomTwitch {
	rand.Seed(time.Now().UnixNano())

	accessToken, err := GetAccessToken()
	if err != nil {
		panic(err)
	}

	streamers, err := GetStreamers(accessToken)
	if err != nil {
		panic(err)
	}


	return RandomTwitch{
		Streamers: streamers,
	}
}

func (rt *RandomTwitch) GetRandomTwitch() string {
	if len(rt.Streamers) == 0 {
		return ""
	}

	selected := rt.Streamers[rand.Intn(len(rt.Streamers))]

	return fmt.Sprintf("%s%s", twichUrl, selected.Username)
}
