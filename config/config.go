package config

import (
	"log"
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

func MustUserID() string {
	cfg := Load()
	if cfg.UserID == "" {
		log.Fatal("TWITCH_USER_ID is not set. Add it to your .env file.")
	}
	return cfg.UserID
}

func ResetCache() {
	loaded = nil
}

func MustUserLogin() string {
	cfg := Load()
	if cfg.UserLogin == "" {
		log.Fatal("TWITCH_USER_LOGIN is not set. Add it to your .env file.")
	}
	return cfg.UserLogin
}
