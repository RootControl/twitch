package api

import (
	"fmt"

	"github.com/RootControl/twitch/config"
	"github.com/RootControl/twitch/entities"
)

type UserResponse struct {
	Data []entities.User `json:"data"`
}

func GetUser(login string) (*UserResponse, error) {
	return getUser(login, nil)
}

func getUser(login string, e Executor) (*UserResponse, error) {
	resp, err := fetch[UserResponse](e, "users", Q("login", login))
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAuthenticatedUser returns the user that owns the token the Twitch CLI is
// configured with. The users endpoint defaults to the token's own user when no
// login or id is supplied.
func GetAuthenticatedUser() (*UserResponse, error) {
	return getAuthenticatedUser(nil)
}

func getAuthenticatedUser(e Executor) (*UserResponse, error) {
	resp, err := fetch[UserResponse](e, "users")
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// cachedUserID memoises the resolved user id for the lifetime of the process.
var cachedUserID string

// resolveUserID determines the id of the current user, in order of preference:
//
//  1. TWITCH_USER_ID, if set
//  2. a lookup of TWITCH_USER_LOGIN, if set
//  3. the user that owns the Twitch CLI's token
//
// This replaces the previous hard requirement that TWITCH_USER_ID be present
// in .env before any followed-streams call would work.
func resolveUserID(e Executor) (string, error) {
	if cachedUserID != "" {
		return cachedUserID, nil
	}

	cfg := config.Load()
	if cfg.UserID != "" {
		cachedUserID = cfg.UserID
		return cachedUserID, nil
	}

	if cfg.UserLogin != "" {
		resp, err := getUser(cfg.UserLogin, e)
		if err != nil {
			return "", fmt.Errorf("looking up TWITCH_USER_LOGIN %q: %w", cfg.UserLogin, err)
		}
		if len(resp.Data) == 0 {
			return "", fmt.Errorf("TWITCH_USER_LOGIN %q does not match a Twitch user", cfg.UserLogin)
		}
		cachedUserID = resp.Data[0].ID
		return cachedUserID, nil
	}

	resp, err := getAuthenticatedUser(e)
	if err != nil {
		return "", fmt.Errorf("determining your user id: %w (set TWITCH_USER_ID in .env to skip this lookup)", err)
	}
	if len(resp.Data) == 0 {
		return "", fmt.Errorf("could not determine your user id; run `twitch token -u` or set TWITCH_USER_ID in .env")
	}
	cachedUserID = resp.Data[0].ID
	return cachedUserID, nil
}

// ResetUserIDCache clears the memoised user id. Used by tests.
func ResetUserIDCache() {
	cachedUserID = ""
}
