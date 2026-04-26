package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/RootControl/twitch/config"
)

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestRootHelp(t *testing.T) {
	out, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"streams", "followed", "user", "categories", "follow", "unfollow", "watch"} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing command %q", want)
		}
	}
}

func TestUserCommandRequiresArg(t *testing.T) {
	_, err := executeCommand("user")
	if err == nil {
		t.Error("user command with no args should return an error")
	}
}

func TestCategoriesCommandRequiresArg(t *testing.T) {
	_, err := executeCommand("categories")
	if err == nil {
		t.Error("categories command with no args should return an error")
	}
}

func TestFollowCommandRequiresArg(t *testing.T) {
	_, err := executeCommand("follow")
	if err == nil {
		t.Error("follow command with no args should return an error")
	}
}

func TestUnfollowCommandRequiresArg(t *testing.T) {
	_, err := executeCommand("unfollow")
	if err == nil {
		t.Error("unfollow command with no args should return an error")
	}
}

func TestWatchCommandHasIntervalFlag(t *testing.T) {
	out, err := executeCommand("watch", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "interval") {
		t.Errorf("watch --help missing --interval flag, got: %q", out)
	}
}

func TestFollowedCommandRegistered(t *testing.T) {
	os.Setenv("TWITCH_USER_ID", "123")
	config.ResetCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		config.ResetCache()
	})

	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "followed" {
			found = true
			break
		}
	}
	if !found {
		t.Error("'followed' command not registered on rootCmd")
	}
}

func TestAllCommandsRegistered(t *testing.T) {
	want := []string{"streams", "followed", "user", "categories", "follow", "unfollow", "watch"}
	registered := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		registered[c.Name()] = true
	}
	for _, name := range want {
		if !registered[name] {
			t.Errorf("command %q not registered", name)
		}
	}
}
