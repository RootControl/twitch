package api

const FOLLOWED_CHANNELS = "channels/followed"

func FollowChannel(broadcasterID string) error {
	return followChannel(broadcasterID, nil)
}

func followChannel(broadcasterID string, e Executor) error {
	userID, err := resolveUserID(e)
	if err != nil {
		return err
	}
	return mutate(e, "post", FOLLOWED_CHANNELS,
		Q("broadcaster_id", broadcasterID),
		Q("user_id", userID),
	)
}

func UnfollowChannel(broadcasterID string) error {
	return unfollowChannel(broadcasterID, nil)
}

func unfollowChannel(broadcasterID string, e Executor) error {
	userID, err := resolveUserID(e)
	if err != nil {
		return err
	}
	return mutate(e, "delete", FOLLOWED_CHANNELS,
		Q("broadcaster_id", broadcasterID),
		Q("user_id", userID),
	)
}
