package main

import "github.com/RootControl/twitch/internal"

func main() {
	randomTwitch := internal.NewRandomTwitch()
	url := randomTwitch.GetRandomTwitch()
	err := internal.OpenBrowser(url)

	if err != nil {
		panic(err)
	}
}
