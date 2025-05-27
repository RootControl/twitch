package main

import "github.com/RootControl/twitch/api"

func main() {
	// streamers := api.GetLiveStreams()
	//
	// for _, stream := range streamers.Data {
	// 	println(stream.GetMainInfo())
	// }

	// categories := api.GetCategories("software")
	// for _, category := range categories.Data {
	// 	println(category.ToString())
	// }

	user := api.GetUser("loganeisenhorn")

	println(user.Data[0].ToString())
}
