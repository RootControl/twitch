package api

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/RootControl/twitch/config"
)

// fakeExecutor captures the last argv and returns a fixed payload.
type fakeExecutor struct {
	payload string
	last    []string
	calls   [][]string
	err     error
}

func (f *fakeExecutor) run(args []string) (*bytes.Buffer, error) {
	f.last = args
	f.calls = append(f.calls, args)
	return bytes.NewBufferString(f.payload), f.err
}

// lastCmd renders the captured argv for substring assertions.
func (f *fakeExecutor) lastCmd() string {
	return strings.Join(f.last, " ")
}

// setUserID pins the configured user id for the duration of a test.
func setUserID(t *testing.T, id string) {
	t.Helper()
	os.Setenv("TWITCH_USER_ID", id)
	config.ResetCache()
	ResetUserIDCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_ID")
		config.ResetCache()
		ResetUserIDCache()
	})
}

// ---- Request rendering ----

func TestRequestToString(t *testing.T) {
	r := NewApiRequest()
	r.Method = "get"
	r.Template = "streams"
	r.Params = []Param{Q("type", "live")}
	got := r.ToString()
	for _, want := range []string{"twitch api", "get", "streams", "-q type=live"} {
		if !strings.Contains(got, want) {
			t.Errorf("ToString() missing %q: %q", want, got)
		}
	}
}

// Arguments must reach the CLI as discrete argv entries, never as a shell
// string, so values with spaces or metacharacters stay intact.
func TestArgsQuoteUnsafeValues(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	if _, err := getCategories("just chatting; id", 50, fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var found bool
	for _, a := range fe.last {
		if a == "query=just chatting; id" {
			found = true
		}
	}
	if !found {
		t.Errorf("query value was split or mangled; argv = %#v", fe.last)
	}
	if len(fe.last) != 7 {
		t.Errorf("expected 7 argv entries (api get tmpl + 2 params), got %d: %#v", len(fe.last), fe.last)
	}
}

func TestRequestMethodsSetFields(t *testing.T) {
	fe := &fakeExecutor{payload: "{}"}
	r := newRequestWithExecutor(fe.run)

	r.Get("streams", Q("first", "1"))
	if r.Method != "get" {
		t.Errorf("Method = %q after Get(), want 'get'", r.Method)
	}
	if r.Template != "streams" {
		t.Errorf("Template = %q after Get(), want 'streams'", r.Template)
	}

	r.Post("channels/followed", Q("broadcaster_id", "1"))
	if r.Method != "post" {
		t.Errorf("Method = %q after Post(), want 'post'", r.Method)
	}

	r.Delete("channels/followed", Q("broadcaster_id", "1"))
	if r.Method != "delete" {
		t.Errorf("Method = %q after Delete(), want 'delete'", r.Method)
	}
}

// ---- error handling ----

func TestFetchSurfacesExecutorError(t *testing.T) {
	fe := &fakeExecutor{payload: "", err: os.ErrPermission}
	if _, err := getLiveStreams(0, fe.run); err == nil {
		t.Error("expected an error when the executor fails")
	}
}

func TestFetchSurfacesAPIErrorPayload(t *testing.T) {
	fe := &fakeExecutor{payload: `{"error":"Unauthorized","status":401,"message":"Invalid OAuth token"}`}
	_, err := getLiveStreams(0, fe.run)
	if err == nil {
		t.Fatal("expected an error for a 401 payload")
	}
	if !strings.Contains(err.Error(), "Invalid OAuth token") {
		t.Errorf("error should quote the API message, got: %v", err)
	}
	if !strings.Contains(err.Error(), "twitch token") {
		t.Errorf("a 401 should suggest re-authenticating, got: %v", err)
	}
}

// The CLI prints the API error envelope on stdout and exits non-zero. The
// envelope explains the failure; the bare exit status does not.
func TestFetchPrefersAPIErrorOverExitStatus(t *testing.T) {
	fe := &fakeExecutor{
		payload: `{"data":[],"error":"Bad Request","status":400,"message":"Bad Identifiers."}`,
		err:     os.ErrPermission,
	}
	_, err := getUser("x; id", fe.run)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "Bad Identifiers.") {
		t.Errorf("error should quote the API message, got: %v", err)
	}
}

func TestFetchIgnoresSuccessfulStatusField(t *testing.T) {
	// A 2xx status in the payload must not be treated as a failure.
	fe := &fakeExecutor{payload: `{"data":[],"status":200}`}
	if _, err := getLiveStreams(0, fe.run); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFetchRejectsEmptyBody(t *testing.T) {
	fe := &fakeExecutor{payload: "  "}
	if _, err := getLiveStreams(0, fe.run); err == nil {
		t.Error("expected an error for an empty response body")
	}
}

func TestFetchRejectsInvalidJSON(t *testing.T) {
	fe := &fakeExecutor{payload: "not json"}
	if _, err := getLiveStreams(0, fe.run); err == nil {
		t.Error("expected an error for a malformed response body")
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
	resp, err := getLiveStreams(0, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	if !strings.Contains(fe.lastCmd(), "streams") {
		t.Errorf("command missing 'streams': %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "type=live") {
		t.Errorf("command missing type=live filter: %q", fe.lastCmd())
	}
}

func TestGetLiveStreamsLimitIsClamped(t *testing.T) {
	cases := []struct {
		limit int
		want  string
	}{
		{0, "first=30"},
		{-5, "first=30"},
		{10, "first=10"},
		{500, "first=100"},
	}
	for _, tc := range cases {
		fe := &fakeExecutor{payload: streamsPayload}
		if _, err := getLiveStreams(tc.limit, fe.run); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(fe.lastCmd(), tc.want) {
			t.Errorf("limit %d: command missing %q: %q", tc.limit, tc.want, fe.lastCmd())
		}
	}
}

func TestGetFollowedStreams(t *testing.T) {
	setUserID(t, "999")

	fe := &fakeExecutor{payload: streamsPayload}
	resp, err := getFollowedStreams(0, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(resp.Data))
	}
	if !strings.Contains(fe.lastCmd(), "streams/followed") {
		t.Errorf("command missing 'streams/followed': %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "user_id=999") {
		t.Errorf("command missing user_id: %q", fe.lastCmd())
	}
}

func TestGetFollowedStreamsEmptyData(t *testing.T) {
	setUserID(t, "999")

	fe := &fakeExecutor{payload: `{"data":[]}`}
	resp, err := getFollowedStreams(0, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 streams, got %d", len(resp.Data))
	}
}

func TestGetStreamsByGamePrefersExactMatch(t *testing.T) {
	// search/categories ranks a fuzzy match first; the exact name must win.
	fe := &fakeExecutor{payload: `{"data":[
		{"id":"1","name":"Rust Console Edition"},
		{"id":"2","name":"Rust"}
	]}`}

	if _, err := getStreamsByGame("rust", 0, fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fe.lastCmd(), "game_id=2") {
		t.Errorf("expected exact match game_id=2, got: %q", fe.lastCmd())
	}
}

func TestGetStreamsByGameFallsBackToFirstResult(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[{"id":"7","name":"Software and Game Dev"}]}`}
	if _, err := getStreamsByGame("software", 0, fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fe.lastCmd(), "game_id=7") {
		t.Errorf("expected game_id=7, got: %q", fe.lastCmd())
	}
}

func TestGetStreamsByGameUnknownCategory(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	if _, err := getStreamsByGame("zzzz", 0, fe.run); err == nil {
		t.Error("expected an error when no category matches")
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
	resp, err := getUser("testuser", fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	if !strings.Contains(fe.lastCmd(), "login=testuser") {
		t.Errorf("command missing login param: %q", fe.lastCmd())
	}
}

func TestGetUserNotFound(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	resp, err := getUser("nobody", fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 users, got %d", len(resp.Data))
	}
}

// ---- user id resolution ----

func TestResolveUserIDPrefersEnvID(t *testing.T) {
	setUserID(t, "555")

	fe := &fakeExecutor{payload: userPayload}
	id, err := resolveUserID(fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "555" {
		t.Errorf("id = %q, want %q", id, "555")
	}
	if len(fe.calls) != 0 {
		t.Errorf("expected no API call when TWITCH_USER_ID is set, got %#v", fe.calls)
	}
}

func TestResolveUserIDFromLogin(t *testing.T) {
	os.Unsetenv("TWITCH_USER_ID")
	os.Setenv("TWITCH_USER_LOGIN", "testuser")
	config.ResetCache()
	ResetUserIDCache()
	t.Cleanup(func() {
		os.Unsetenv("TWITCH_USER_LOGIN")
		config.ResetCache()
		ResetUserIDCache()
	})

	fe := &fakeExecutor{payload: userPayload}
	id, err := resolveUserID(fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "42" {
		t.Errorf("id = %q, want %q", id, "42")
	}
	if !strings.Contains(fe.lastCmd(), "login=testuser") {
		t.Errorf("expected a users lookup by login, got: %q", fe.lastCmd())
	}
}

func TestResolveUserIDFallsBackToAuthenticatedUser(t *testing.T) {
	os.Unsetenv("TWITCH_USER_ID")
	os.Unsetenv("TWITCH_USER_LOGIN")
	config.ResetCache()
	ResetUserIDCache()
	t.Cleanup(func() {
		config.ResetCache()
		ResetUserIDCache()
	})

	fe := &fakeExecutor{payload: userPayload}
	id, err := resolveUserID(fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "42" {
		t.Errorf("id = %q, want %q", id, "42")
	}
	if strings.Contains(fe.lastCmd(), "login=") {
		t.Errorf("authenticated lookup should send no login param, got: %q", fe.lastCmd())
	}
}

func TestResolveUserIDIsCached(t *testing.T) {
	os.Unsetenv("TWITCH_USER_ID")
	os.Unsetenv("TWITCH_USER_LOGIN")
	config.ResetCache()
	ResetUserIDCache()
	t.Cleanup(func() {
		config.ResetCache()
		ResetUserIDCache()
	})

	fe := &fakeExecutor{payload: userPayload}
	if _, err := resolveUserID(fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := resolveUserID(fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fe.calls) != 1 {
		t.Errorf("expected the id to be resolved once, got %d calls", len(fe.calls))
	}
}

func TestResolveUserIDErrorsWhenLookupFails(t *testing.T) {
	os.Unsetenv("TWITCH_USER_ID")
	os.Unsetenv("TWITCH_USER_LOGIN")
	config.ResetCache()
	ResetUserIDCache()
	t.Cleanup(func() {
		config.ResetCache()
		ResetUserIDCache()
	})

	fe := &fakeExecutor{payload: `{"data":[]}`}
	if _, err := resolveUserID(fe.run); err == nil {
		t.Error("expected an error when the user id cannot be determined")
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
	resp, err := getCategories("just", 50, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 category, got %d", len(resp.Data))
	}
	c := resp.Data[0]
	if c.Name != "Just Chatting" {
		t.Errorf("Name = %q, want %q", c.Name, "Just Chatting")
	}
	if !strings.Contains(fe.lastCmd(), "search/categories") {
		t.Errorf("command missing 'search/categories': %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "query=just") {
		t.Errorf("command missing query param: %q", fe.lastCmd())
	}
}

func TestGetCategoriesEmpty(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	resp, err := getCategories("zzz", 0, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 categories, got %d", len(resp.Data))
	}
}

func TestGetTopGames(t *testing.T) {
	fe := &fakeExecutor{payload: categoriesPayload}
	resp, err := getTopGames(5, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 category, got %d", len(resp.Data))
	}
	if !strings.Contains(fe.lastCmd(), "games/top") {
		t.Errorf("command missing 'games/top': %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "first=5") {
		t.Errorf("command missing limit: %q", fe.lastCmd())
	}
}

// ---- channel search ----

func TestSearchChannels(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[
		{"id":"1","broadcaster_login":"someone","display_name":"Someone","is_live":true,"game_name":"Go","title":"hi"}
	]}`}

	resp, err := searchChannels("some", false, 20, fe.run)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(resp.Data))
	}
	if resp.Data[0].BroadcasterName != "Someone" {
		t.Errorf("BroadcasterName = %q, want %q", resp.Data[0].BroadcasterName, "Someone")
	}
	if !strings.Contains(fe.lastCmd(), "search/channels") {
		t.Errorf("command missing 'search/channels': %q", fe.lastCmd())
	}
	if strings.Contains(fe.lastCmd(), "live_only") {
		t.Errorf("live_only should be omitted when not requested: %q", fe.lastCmd())
	}
}

func TestSearchChannelsLiveOnly(t *testing.T) {
	fe := &fakeExecutor{payload: `{"data":[]}`}
	if _, err := searchChannels("some", true, 20, fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fe.lastCmd(), "live_only=true") {
		t.Errorf("command missing live_only: %q", fe.lastCmd())
	}
}

// ---- follows ----

func TestFollowChannel(t *testing.T) {
	setUserID(t, "111")

	fe := &fakeExecutor{payload: "{}"}
	if err := followChannel("999", fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(fe.lastCmd(), "post") {
		t.Errorf("follow should use POST, got: %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "broadcaster_id=999") {
		t.Errorf("command missing broadcaster_id: %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "user_id=111") {
		t.Errorf("command missing user_id: %q", fe.lastCmd())
	}
}

func TestUnfollowChannel(t *testing.T) {
	setUserID(t, "111")

	fe := &fakeExecutor{payload: "{}"}
	if err := unfollowChannel("999", fe.run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(fe.lastCmd(), "delete") {
		t.Errorf("unfollow should use DELETE, got: %q", fe.lastCmd())
	}
	if !strings.Contains(fe.lastCmd(), "broadcaster_id=999") {
		t.Errorf("command missing broadcaster_id: %q", fe.lastCmd())
	}
}

func TestFollowChannelSurfacesAPIError(t *testing.T) {
	setUserID(t, "111")

	fe := &fakeExecutor{payload: `{"error":"Unauthorized","status":401,"message":"Missing scope"}`}
	if err := followChannel("999", fe.run); err == nil {
		t.Error("expected an error for a 401 response")
	}
}

// ---- orShell ----

func TestOrShellReturnsProvidedExecutor(t *testing.T) {
	fe := &fakeExecutor{payload: ""}
	got := orShell(fe.run)
	got([]string{"echo", "test"})
	if fe.lastCmd() != "echo test" {
		t.Errorf("orShell returned wrong executor, last = %q", fe.lastCmd())
	}
}

func TestOrShellFallsBackToShell(t *testing.T) {
	got := orShell(nil)
	if got == nil {
		t.Error("orShell(nil) returned nil")
	}
}
