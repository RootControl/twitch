package player

import (
	"fmt"
	"os/exec"
)

type Player struct {
	Name string
	bin  string
	args func(url string) []string
}

var supported = []Player{
	{
		Name: "streamlink",
		bin:  "streamlink",
		args: func(url string) []string { return []string{url, "best"} },
	},
	{
		Name: "mpv",
		bin:  "mpv",
		args: func(url string) []string { return []string{url} },
	},
}

// Available returns installed players from the supported list.
func Available() []Player {
	var out []Player
	for _, p := range supported {
		if _, err := exec.LookPath(p.bin); err == nil {
			out = append(out, p)
		}
	}
	return out
}

func (p Player) Open(url string) error {
	cmd := exec.Command(p.bin, p.args(url)...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not launch %s: %w", p.Name, err)
	}
	return nil
}
