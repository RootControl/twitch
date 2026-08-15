package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs the twitch CLI with the given argv and returns its stdout.
type Executor func(args []string) (*bytes.Buffer, error)

// ErrCLIMissing is returned when the Twitch CLI is not installed or not on PATH.
var ErrCLIMissing = errors.New("the Twitch CLI is not installed or not on PATH (see https://dev.twitch.tv/docs/cli/)")

// shellExecutor runs the twitch binary directly. Arguments are passed as argv
// rather than through a shell, so values containing spaces or shell
// metacharacters are handed to the CLI verbatim.
func shellExecutor(args []string) (*bytes.Buffer, error) {
	cmd := exec.Command("twitch", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return &out, ErrCLIMissing
		}
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return &out, fmt.Errorf("twitch %s: %w: %s", strings.Join(args, " "), err, msg)
		}
		return &out, fmt.Errorf("twitch %s: %w", strings.Join(args, " "), err)
	}
	return &out, nil
}

// Param is a single -q query parameter passed to the Twitch CLI.
type Param struct {
	Key   string
	Value string
}

// Q builds a query parameter.
func Q(key, value string) Param {
	return Param{Key: key, Value: value}
}

type Request struct {
	command  string
	Method   string
	Template string
	Params   []Param
	exec     Executor
}

func NewApiRequest() *Request {
	return &Request{
		command: "twitch",
		Params:  make([]Param, 0),
		exec:    shellExecutor,
	}
}

func newRequestWithExecutor(e Executor) *Request {
	r := NewApiRequest()
	r.exec = e
	return r
}

// Args returns the argv passed to the twitch binary (excluding the binary itself).
func (r *Request) Args() []string {
	args := make([]string, 0, 3+len(r.Params)*2)
	args = append(args, "api", r.Method, r.Template)
	for _, p := range r.Params {
		args = append(args, "-q", p.Key+"="+p.Value)
	}
	return args
}

// ToString renders the command for display and debugging only. It is never
// handed to a shell.
func (r *Request) ToString() string {
	return strings.TrimSpace(r.command + " " + strings.Join(r.Args(), " "))
}

func (r *Request) run() (*bytes.Buffer, error) {
	return r.exec(r.Args())
}

func (r *Request) do(method, template string, params ...Param) (*bytes.Buffer, error) {
	r.Method = method
	r.Template = template
	r.Params = params
	return r.run()
}

func (r *Request) Get(template string, params ...Param) (*bytes.Buffer, error) {
	return r.do("get", template, params...)
}

func (r *Request) Post(template string, params ...Param) (*bytes.Buffer, error) {
	return r.do("post", template, params...)
}

func (r *Request) Delete(template string, params ...Param) (*bytes.Buffer, error) {
	return r.do("delete", template, params...)
}

// apiError models the error payload the Helix API returns for non-2xx
// responses. The CLI prints it on stdout with a zero exit code, so without
// this check an expired token looks identical to an empty result set.
type apiError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func checkAPIError(body []byte) error {
	var e apiError
	if err := json.Unmarshal(body, &e); err != nil {
		return nil // not an error envelope; let the caller decode normally
	}
	if e.Status < 400 {
		return nil
	}
	msg := e.Message
	if msg == "" {
		msg = e.Error
	}
	if e.Status == 401 {
		return fmt.Errorf("twitch API %d: %s (run `twitch token` to re-authenticate)", e.Status, msg)
	}
	return fmt.Errorf("twitch API %d: %s", e.Status, msg)
}

// body safely extracts the bytes captured from a run, which may be nil when
// the command could not be started at all.
func body(buf *bytes.Buffer) []byte {
	if buf == nil {
		return nil
	}
	return buf.Bytes()
}

// fetch runs a GET and decodes the JSON body into T.
func fetch[T any](e Executor, template string, params ...Param) (T, error) {
	var result T

	request := newRequestWithExecutor(orShell(e))
	buf, execErr := request.Get(template, params...)
	payload := body(buf)

	// The CLI writes the API's error envelope to stdout and exits non-zero,
	// so check the payload before the exit status: it carries the message
	// that actually explains the failure.
	if err := checkAPIError(payload); err != nil {
		return result, err
	}
	if execErr != nil {
		return result, execErr
	}

	if len(bytes.TrimSpace(payload)) == 0 {
		return result, fmt.Errorf("empty response from %q", template)
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		return result, fmt.Errorf("parsing response from %q: %w", template, err)
	}
	return result, nil
}

// mutate runs a POST/DELETE that returns no meaningful body.
func mutate(e Executor, method, template string, params ...Param) error {
	request := newRequestWithExecutor(orShell(e))
	buf, execErr := request.do(method, template, params...)

	if err := checkAPIError(body(buf)); err != nil {
		return err
	}
	return execErr
}

func orShell(e Executor) Executor {
	if e != nil {
		return e
	}
	return shellExecutor
}
