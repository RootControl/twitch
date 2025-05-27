package main

import "github.com/RootControl/twitch/api"

func main() {
	streamers := api.GetLiveStreams()

	for _, stream := range streamers.Data {
		println(stream.GetMainInfo())
	}
}
