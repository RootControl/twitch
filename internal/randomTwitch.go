package internal

import (
	"fmt"
	"time"
)

type RandomTwitch struct {
	urls []string
}

const twichUrl = "https://www.twitch.tv/"

func NewRandomTwitch() RandomTwitch {
	return RandomTwitch{
		urls: []string{
			"nyybeats",
			}
	}
}

func (rt *RandomTwitch) GetRandomTwitch() string {
	rand.seed(time.Now().UnixNano())
	randomTwitch := rt.urls[rand.Intn(len(rt.urls))]

	return fmt.Sprintf("%s%s", twichUrl, randomTwitch)
}
