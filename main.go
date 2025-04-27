package main

import "github.com/RootControl/twitch/internal"

func main() {
	randomTwitch := internal.NewRandomTwitch()
	url := randomTwitch.GetRandomTwitch()
	internal.OpenBrowser(url)
}
