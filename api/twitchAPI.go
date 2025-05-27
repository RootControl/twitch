package api

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Request struct {
	command  string
	Method   string
	Template string
	Flags    []string
}

func NewApiRequest() *Request {
	return &Request{
		command:  "twitch api",
		Method:   "",
		Template: "",
		Flags:    make([]string, 0),
	}
}

func (r *Request) ToString() string {
	return fmt.Sprintf("%s %s %s %s", r.command, r.Method, r.Template, strings.Join(r.Flags, " "))
}

func (r *Request) Get(template string, flags ...string) *bytes.Buffer {
	r.Method = "get"
	r.Template = template
	r.Flags = flags

	cmd := exec.Command("sh", "-c", r.ToString())

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running twitch CLI: %v\nOutput: %s", err, out.String())

	}

	return &out
}
