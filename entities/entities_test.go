package entities

import (
	"strings"
	"testing"
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

func TestNewStream(t *testing.T) {
	s := NewStream()
	if s == nil {
		t.Fatal("NewStream() returned nil")
	}
}

func TestNewUser(t *testing.T) {
	u := NewUser()
	if u == nil {
		t.Fatal("NewUser() returned nil")
	}
}

func TestNewCategory(t *testing.T) {
	c := NewCategory()
	if c == nil {
		t.Fatal("NewCategory() returned nil")
	}
}
