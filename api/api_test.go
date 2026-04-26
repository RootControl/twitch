package api

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/RootControl/twitch/config"
)

// fakeExecutor captures the last command string and returns a fixed payload.
type fakeExecutor struct {
	payload string
	last    string
	err     error
}

func (f *fakeExecutor) run(cmd string) (*bytes.Buffer, error) {
	f.last = cmd
	return bytes.NewBufferString(f.payload), f.err
}

// ---- Request.ToString ----

func TestRequestToString(t *testing.T) {
	r := NewApiRequest()
	r.Method = "get"
	r.Template = "streams"
	r.Flags = []string{"-q type=live"}
	got := r.ToString()
	if !strings.Contains(got, "twitch api") {
		t.Errorf("ToString() missing 'twitch api': %q", got)
	}
	if !strings.Contains(got, "get") {
		t.Errorf("ToString() missing method 'get': %q", got)
	}
	if !strings.Contains(got, "streams") {
		t.Errorf("ToString() missing template 'streams': %q", got)
	}
	if !strings.Contains(got, "-q type=live") {
		t.Errorf("ToString() missing flag: %q", got)
	}
}

func TestRequestMethodsSetFields(t *testing.T) {
	fe := &fakeExecutor{payload: "{}"}
	r := newRequestWithExecutor(fe.run)

	r.Get("streams", "-q first=1")
	if r.Method != "get" {
		t.Errorf("Method = %q after Get(), want 'get'", r.Method)
	}
	if r.Template != "streams" {
		t.Errorf("Template = %q after Get(), want 'streams'", r.Template)
	}

	r.Post("channels/followed", "-q broadcaster_id=1")
	if r.Method != "post" {
		t.Errorf("Method = %q after Post(), want 'post'", r.Method)
	}

	r.Delete("channels/followed", "-q broadcaster_id=1")
	if r.Method != "delete" {
		t.Errorf("Method = %q after Delete(), want 'delete'", r.Method)
	}
}

// ---- streamers ----

var streamsPayload = `{
  "data": [
    {
      "id": "1",
      "user_login": "streamer1",
      "user_name": "Streamer1",
      "game_name": "Go",
      "title": "Coding live",
      "viewer_count": 100
    }
  ]
}`

func TestGetLiveStreams(t *testing.T) {
	fe := &fakeExecutor{payload: streamsPayload}
	resp := getLiveStreams(fe.run)

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(resp.Data))
	}
	s := resp.Data[0]
	if s.Username != "Streamer1" {
		t.Errorf("Username = %q, want %q", s.Username, "Streamer1")
	}
	if s.GameName != "Go" {
		t.Errorf("GameName = %q, want %q", s.GameName, "Go")
	}
	if !strings.Contains(fe.last, "streams") {
		t.Errorf("command missing 'streams': %q", fe.last)
	}
	if !strings.Contains(fe.last, "type=live") {
		t.Errorf("command missing type=live filter: %q", fe.last)
	}
}

func TestGetFollowedStreams(t *testing.T) {
	os.Setenv("TWITCH_USER_ID", "999")
	config.ResetCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		config.ResetCache()
	})

	fe := &fakeExecutor{payload: streamsPayload}
	resp := getFollowedStreams(fe.run)

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(resp.Data))
	}
	if !strings.Contains(fe.last, "streams/followed") {
		t.Errorf("command missing 'streams/followed': %q", fe.last)
	}
	if !strings.Contains(fe.last, "user_id=999") {
		t.Errorf("command missing user_id: %q", fe.last)
	}
}

func TestGetFollowedStreamsEmptyData(t *testing.T) {
	os.Setenv("TWITCH_USER_ID", "999")
	config.ResetCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		config.ResetCache()
	})

	fe := &fakeExecutor{payload: `{"data":[]}`}
	resp := getFollowedStreams(fe.run)
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 streams, got %d", len(resp.Data))
	}
}

// ---- users ----

var userPayload = `{
  "data": [
    {
      "id": "42",
      "login": "testuser",
      "display_name": "TestUser"
    }
  ]
}`

func TestGetUser(t *testing.T) {
	fe := &fakeExecutor{payload: userPayload}
	resp := getUser("testuser", fe.run)

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 user, got %d", len(resp.Data))
	}
	u := resp.Data[0]
	if u.Login != "testuser" {
		t.Errorf("Login = %q, want %q", u.Login, "testuser")
	}
	if u.ID != "42" {
		t.Errorf("ID = %q, want %q", u.ID, "42")
	}
	if !strings.Contains(fe.last, "login=testuser") {
		t.Errorf("command missing login param: %q", fe.last)
	}
}

func TestGetUserNotFound(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	resp := getUser("nobody", fe.run)
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 users, got %d", len(resp.Data))
	}
}

// ---- categories ----

var categoriesPayload = `{
  "data": [
    {"id": "509658", "name": "Just Chatting", "box_art_url": "https://example.com/art.jpg"}
  ],
  "pagination": {}
}`

func TestGetCategories(t *testing.T) {
	fe := &fakeExecutor{payload: categoriesPayload}
	resp := getCategories("just", fe.run)

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 category, got %d", len(resp.Data))
	}
	c := resp.Data[0]
	if c.Name != "Just Chatting" {
		t.Errorf("Name = %q, want %q", c.Name, "Just Chatting")
	}
	if !strings.Contains(fe.last, "search/categories") {
		t.Errorf("command missing 'search/categories': %q", fe.last)
	}
	if !strings.Contains(fe.last, "query=just") {
		t.Errorf("command missing query param: %q", fe.last)
	}
}

func TestGetCategoriesEmpty(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	resp := getCategories("zzz", fe.run)
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 categories, got %d", len(resp.Data))
	}
}

// ---- follows ----

func TestFollowChannel(t *testing.T) {
	os.Setenv("TWITCH_USER_ID", "111")
	config.ResetCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		config.ResetCache()
	})

	fe := &fakeExecutor{payload: ""}
	followChannel("999", fe.run)

	if !strings.Contains(fe.last, "post") {
		t.Errorf("follow should use POST, got: %q", fe.last)
	}
	if !strings.Contains(fe.last, "broadcaster_id=999") {
		t.Errorf("command missing broadcaster_id: %q", fe.last)
	}
	if !strings.Contains(fe.last, "user_id=111") {
		t.Errorf("command missing user_id: %q", fe.last)
	}
}

func TestUnfollowChannel(t *testing.T) {
	os.Setenv("TWITCH_USER_ID", "111")
	config.ResetCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		config.ResetCache()
	})

	fe := &fakeExecutor{payload: ""}
	unfollowChannel("999", fe.run)

	if !strings.Contains(fe.last, "delete") {
		t.Errorf("unfollow should use DELETE, got: %q", fe.last)
	}
	if !strings.Contains(fe.last, "broadcaster_id=999") {
		t.Errorf("command missing broadcaster_id: %q", fe.last)
	}
}

// ---- orShell ----

func TestOrShellReturnsProvidedExecutor(t *testing.T) {
	fe := &fakeExecutor{payload: ""}
	got := orShell(fe.run)
	// call it and check fe.last was set (proves it's the same func)
	got("echo test")
	if fe.last != "echo test" {
		t.Errorf("orShell returned wrong executor")
	}
}

func TestOrShellFallsBackToShell(t *testing.T) {
	got := orShell(nil)
	if got == nil {
		t.Error("orShell(nil) returned nil")
	}
}
