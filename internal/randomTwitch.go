package internal

import (
	"fmt"
	"math/rand"
	"time"
)

type RandomTwitch struct {
	urls []string
}

const twichUrl = "https://www.twitch.tv/"

func NewRandomTwitch() RandomTwitch {
	return RandomTwitch{
		urls: []string{
			"xqc",
			"pokimane",
			"shroud",
			"lirik",
			"ninja",
		},
	}
}

func (rt *RandomTwitch) GetRandomTwitch() string {
	rand.Seed(time.Now().UnixNano())
	randomTwitch := rt.urls[rand.Intn(len(rt.urls))]

	return fmt.Sprintf("%s%s", twichUrl, randomTwitch)
}
