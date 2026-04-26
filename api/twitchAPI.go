package api

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Executor runs a shell command and returns its stdout.
type Executor func(command string) (*bytes.Buffer, error)

func shellExecutor(command string) (*bytes.Buffer, error) {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return &out, err
}

type Request struct {
	command  string
	Method   string
	Template string
	Flags    []string
	exec     Executor
}

func NewApiRequest() *Request {
	return &Request{
		command: "twitch api",
		Flags:   make([]string, 0),
		exec:    shellExecutor,
	}
}

func newRequestWithExecutor(e Executor) *Request {
	r := NewApiRequest()
	r.exec = e
	return r
}

func (r *Request) ToString() string {
	return fmt.Sprintf("%s %s %s %s", r.command, r.Method, r.Template, strings.Join(r.Flags, " "))
}

func (r *Request) run() *bytes.Buffer {
	out, err := r.exec(r.ToString())
	if err != nil {
		log.Fatalf("Error running twitch CLI: %v\nOutput: %s", err, out.String())
	}
	return out
}

func (r *Request) Get(template string, flags ...string) *bytes.Buffer {
	r.Method = "get"
	r.Template = template
	r.Flags = flags
	return r.run()
}

func (r *Request) Post(template string, flags ...string) *bytes.Buffer {
	r.Method = "post"
	r.Template = template
	r.Flags = flags
	return r.run()
}

func (r *Request) Delete(template string, flags ...string) *bytes.Buffer {
	r.Method = "delete"
	r.Template = template
	r.Flags = flags
	return r.run()
}
