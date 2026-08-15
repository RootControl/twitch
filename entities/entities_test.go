package entities

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestStreamGetMainInfo(t *testing.T) {
	s := Stream{
		Username:  "streamer1",
		Title:     "Playing Go",
		GameName:  "Software and Game Dev",
		UserLogin: "streamer1",
	}
	info := s.GetMainInfo()
	for _, want := range []string{"streamer1", "Playing Go", "Software and Game Dev", TwitchURL + "streamer1"} {
		if !strings.Contains(info, want) {
			t.Errorf("GetMainInfo() missing %q in %q", want, info)
		}
	}
}

func TestStreamURL(t *testing.T) {
	s := Stream{UserLogin: "someone"}
	if got := s.URL(); got != TwitchURL+"someone" {
		t.Errorf("URL() = %q, want %q", got, TwitchURL+"someone")
	}
}

func TestUptimeSince(t *testing.T) {
	now := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name      string
		startedAt string
		want      string
	}{
		{"two hours", "2026-08-14T10:00:00Z", "2h00m"},
		{"minutes only", "2026-08-14T11:43:00Z", "17m"},
		{"hours and minutes", "2026-08-14T09:46:00Z", "2h14m"},
		{"empty", "", ""},
		{"unparseable", "yesterday", ""},
		{"in the future", "2026-08-14T13:00:00Z", ""},
	}
	for _, tc := range cases {
		if got := uptimeSince(tc.startedAt, now); got != tc.want {
			t.Errorf("%s: uptimeSince(%q) = %q, want %q", tc.name, tc.startedAt, got, tc.want)
		}
	}
}

// Helix returns the field as tag_ids; the previous tags_ids tag silently
// dropped it.
func TestStreamDecodesTagIDs(t *testing.T) {
	var s Stream
	if err := json.Unmarshal([]byte(`{"tag_ids":["abc","def"]}`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.TagsIds) != 2 {
		t.Errorf("TagsIds = %#v, want 2 entries", s.TagsIds)
	}
}

// promptui renders items through text/template and passes them by value. A
// pointer method is not in the method set of a non-addressable value, so with
// a pointer receiver the template cannot resolve these and renders the raw
// struct instead.
func TestStreamMethodsResolveFromTemplate(t *testing.T) {
	s := Stream{
		Username:  "streamer1",
		UserLogin: "streamer1",
		Title:     "Coding",
		StartedAt: "2020-01-01T00:00:00Z",
	}

	for _, expr := range []string{"{{ .Uptime }}", "{{ .URL }}", "{{ .GetMainInfo }}"} {
		tmpl, err := template.New("t").Parse(expr)
		if err != nil {
			t.Fatalf("parsing %s: %v", expr, err)
		}

		var buf bytes.Buffer
		// Pass by value, exactly as promptui does.
		if err := tmpl.Execute(&buf, s); err != nil {
			t.Errorf("executing %s: %v", expr, err)
			continue
		}
		if buf.Len() == 0 {
			t.Errorf("%s rendered empty", expr)
		}
	}
}

func TestUserToString(t *testing.T) {
	u := User{ID: "123", Login: "user1", DisplayName: "User One"}
	out := u.ToString()
	for _, want := range []string{"123", "user1", "User One"} {
		if !strings.Contains(out, want) {
			t.Errorf("ToString() missing %q in %q", want, out)
		}
	}
}

func TestCategoryToString(t *testing.T) {
	c := Category{ID: "509658", Name: "Just Chatting"}
	out := c.ToString()
	for _, want := range []string{"509658", "Just Chatting"} {
		if !strings.Contains(out, want) {
			t.Errorf("ToString() missing %q in %q", want, out)
		}
	}
}

func TestChannelToString(t *testing.T) {
	c := Channel{
		BroadcasterName:     "Someone",
		BroadcasterLogin:    "someone",
		BroadcasterLanguage: "en",
		GameName:            "Go",
		Title:               "Coding",
	}
	out := c.ToString()
	for _, want := range []string{"Someone", "Coding", "Go", "en", TwitchURL + "someone"} {
		if !strings.Contains(out, want) {
			t.Errorf("ToString() missing %q in %q", want, out)
		}
	}
}

func TestSearchedChannelToString(t *testing.T) {
	live := SearchedChannel{BroadcasterName: "Someone", BroadcasterLogin: "someone", IsLive: true, Title: "Coding", GameName: "Go"}
	if !strings.Contains(live.ToString(), "LIVE") {
		t.Errorf("live channel should be marked LIVE: %q", live.ToString())
	}

	offline := SearchedChannel{BroadcasterName: "Someone", BroadcasterLogin: "someone"}
	if !strings.Contains(offline.ToString(), "offline") {
		t.Errorf("offline channel should be marked offline: %q", offline.ToString())
	}
}

// search/channels names its id field "id", not "broadcaster_id".
func TestSearchedChannelDecodesID(t *testing.T) {
	var c SearchedChannel
	if err := json.Unmarshal([]byte(`{"id":"77","display_name":"Someone","is_live":true}`), &c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BroadcasterID != "77" {
		t.Errorf("BroadcasterID = %q, want %q", c.BroadcasterID, "77")
	}
	if c.BroadcasterName != "Someone" {
		t.Errorf("BroadcasterName = %q, want %q", c.BroadcasterName, "Someone")
	}
}

func TestNewStream(t *testing.T) {
	if s := NewStream(); s == nil {
		t.Fatal("NewStream() returned nil")
	}
}

func TestNewUser(t *testing.T) {
	if u := NewUser(); u == nil {
		t.Fatal("NewUser() returned nil")
	}
}

func TestNewCategory(t *testing.T) {
	if c := NewCategory(); c == nil {
		t.Fatal("NewCategory() returned nil")
	}
}
