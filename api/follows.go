package api

import (
	"github.com/RootControl/twitch/config"
)

func FollowChannel(broadcasterID string) {
	followChannel(broadcasterID, nil)
}

func followChannel(broadcasterID string, e Executor) {
	request := newRequestWithExecutor(orShell(e))
	request.Post("channels/followed",
		"-q broadcaster_id="+broadcasterID,
		"-q user_id="+config.MustUserID(),
	)
}

func UnfollowChannel(broadcasterID string) {
	unfollowChannel(broadcasterID, nil)
}

func unfollowChannel(broadcasterID string, e Executor) {
	request := newRequestWithExecutor(orShell(e))
	request.Delete("channels/followed",
		"-q broadcaster_id="+broadcasterID,
		"-q user_id="+config.MustUserID(),
	)
}
