package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	UserID    string
	UserLogin string
}

var loaded *Config

func Load() *Config {
	if loaded != nil {
		return loaded
	}

	_ = godotenv.Load()

	loaded = &Config{
		UserID:    os.Getenv("TWITCH_USER_ID"),
		UserLogin: os.Getenv("TWITCH_USER_LOGIN"),
	}
	return loaded
}

// UserID returns the configured user id, or an error naming the setting that
// is missing. Callers that can fall back to a lookup should read Load()
// directly rather than treating an unset value as fatal.
func UserID() (string, error) {
	cfg := Load()
	if cfg.UserID == "" {
		return "", fmt.Errorf("TWITCH_USER_ID is not set")
	}
	return cfg.UserID, nil
}

// UserLogin returns the configured login, or an error if it is unset.
func UserLogin() (string, error) {
	cfg := Load()
	if cfg.UserLogin == "" {
		return "", fmt.Errorf("TWITCH_USER_LOGIN is not set")
	}
	return cfg.UserLogin, nil
}

func ResetCache() {
	loaded = nil
}
