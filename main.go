package main

import "github.com/RootControl/Twitch/internal"

func main() {
	randomTwitch := internal.NewRandomTwitch()
	url := randomTwitch.GetRandomTwitch()
	err := internal.OpenBrowser(url)

	if err != nil {
		panic(err)
	}
}
