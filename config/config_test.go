package config

import (
	"os"
	"testing"
)

func resetCache() {
	loaded = nil
}

func TestLoadReadsEnvVars(t *testing.T) {
	resetCache()
	os.Setenv("TWITCH_USER_ID", "12345")
	os.Setenv("TWITCH_USER_LOGIN", "testuser")
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		os.Unsetenv("TWITCH_USER_LOGIN")
		resetCache()
	})

	cfg := Load()
	if cfg.UserID != "12345" {
		t.Errorf("UserID = %q, want %q", cfg.UserID, "12345")
	}
	if cfg.UserLogin != "testuser" {
		t.Errorf("UserLogin = %q, want %q", cfg.UserLogin, "testuser")
	}
}

func TestLoadReturnsCachedValue(t *testing.T) {
	resetCache()
	os.Setenv("TWITCH_USER_ID", "first")
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		resetCache()
	})

	cfg1 := Load()
	os.Setenv("TWITCH_USER_ID", "second")
	cfg2 := Load()

	if cfg1 != cfg2 {
		t.Error("Load() should return the same cached pointer on repeated calls")
	}
	if cfg2.UserID != "first" {
		t.Errorf("cached UserID = %q, want %q", cfg2.UserID, "first")
	}
}

func TestLoadEmptyWhenEnvNotSet(t *testing.T) {
	resetCache()
	os.Unsetenv("TWITCH_USER_ID")
	os.Unsetenv("TWITCH_USER_LOGIN")
	t.Cleanup(resetCache)

	cfg := Load()
	if cfg.UserID != "" {
		t.Errorf("UserID should be empty, got %q", cfg.UserID)
	}
	if cfg.UserLogin != "" {
		t.Errorf("UserLogin should be empty, got %q", cfg.UserLogin)
	}
}

func TestUserIDReturnsValue(t *testing.T) {
	resetCache()
	os.Setenv("TWITCH_USER_ID", "99")
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		resetCache()
	})

	id, err := UserID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "99" {
		t.Errorf("UserID() = %q, want %q", id, "99")
	}
}

func TestUserIDErrorsWhenUnset(t *testing.T) {
	resetCache()
	os.Unsetenv("TWITCH_USER_ID")
	t.Cleanup(resetCache)

	if _, err := UserID(); err == nil {
		t.Error("UserID() should return an error when TWITCH_USER_ID is unset")
	}
}

func TestUserLoginReturnsValue(t *testing.T) {
	resetCache()
	os.Setenv("TWITCH_USER_LOGIN", "mylogin")
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_LOGIN")
		resetCache()
	})

	login, err := UserLogin()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if login != "mylogin" {
		t.Errorf("UserLogin() = %q, want %q", login, "mylogin")
	}
}

func TestUserLoginErrorsWhenUnset(t *testing.T) {
	resetCache()
	os.Unsetenv("TWITCH_USER_LOGIN")
	t.Cleanup(resetCache)

	if _, err := UserLogin(); err == nil {
		t.Error("UserLogin() should return an error when TWITCH_USER_LOGIN is unset")
	}
}
