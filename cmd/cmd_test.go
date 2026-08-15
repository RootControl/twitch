package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/RootControl/twitch/config"
	"github.com/RootControl/twitch/entities"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func executeCommand(args ...string) (string, error) {
	// rootCmd is a package-level singleton, so flag state (notably --help)
	// survives between invocations and would leak into the next test.
	resetFlags(rootCmd)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetFlags(c *cobra.Command) {
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			_ = f.Value.Set(f.DefValue)
			f.Changed = false
		}
	})
	for _, sub := range c.Commands() {
		resetFlags(sub)
	}
}

var allCommands = []string{
	"streams", "followed", "user", "channel", "categories",
	"top", "search", "open", "follow", "unfollow", "watch",
}

func TestRootHelp(t *testing.T) {
	out, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range allCommands {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing command %q", want)
		}
	}
}

func TestAllCommandsRegistered(t *testing.T) {
	registered := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		registered[c.Name()] = true
	}
	for _, name := range allCommands {
		if !registered[name] {
			t.Errorf("command %q not registered", name)
		}
	}
}

func TestCommandsRequiringOneArg(t *testing.T) {
	for _, name := range []string{"user", "channel", "categories", "search", "open", "follow", "unfollow"} {
		if _, err := executeCommand(name); err == nil {
			t.Errorf("%s with no args should return an error", name)
		}
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

// A too-small interval must be rejected up front rather than spinning the
// poller in a tight loop against the API.
func TestWatchRejectsShortInterval(t *testing.T) {
	t.Cleanup(func() { watchInterval = 60 })

	if _, err := executeCommand("watch", "--interval", "1"); err == nil {
		t.Error("watch --interval 1 should return an error")
	}
}

func TestLimitFlagsRegistered(t *testing.T) {
	for _, name := range []string{"streams", "followed", "categories", "top", "search"} {
		out, err := executeCommand(name, "--help")
		if err != nil {
			t.Fatalf("%s --help: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "--limit") {
			t.Errorf("%s --help missing --limit flag", name)
		}
	}
}

func TestCompletionRejectsUnknownShell(t *testing.T) {
	if _, err := executeCommand("completion", "cmd.exe"); err == nil {
		t.Error("completion with an unsupported shell should return an error")
	}
}

func TestCompletionAcceptsKnownShell(t *testing.T) {
	out, err := executeCommand("completion", "bash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("completion bash produced no output")
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

// ---- invoked name ----

func TestExecutableName(t *testing.T) {
	orig := os.Args
	t.Cleanup(func() { os.Args = orig })

	cases := []struct {
		argv0 string
		want  string
	}{
		{"ttv", "ttv"},
		{"./ttv", "ttv"},
		{"/usr/local/bin/ttv", "ttv"},
		{"/usr/local/bin/ttv.exe", "ttv"},
		{"twitch", "twitch"},
		{"", defaultName},
		{"/", defaultName},
	}
	for _, tc := range cases {
		os.Args = []string{tc.argv0}
		if got := executableName(); got != tc.want {
			t.Errorf("executableName() with argv[0]=%q = %q, want %q", tc.argv0, got, tc.want)
		}
	}
}

func TestExecutableNameWithoutArgs(t *testing.T) {
	orig := os.Args
	t.Cleanup(func() { os.Args = orig })

	os.Args = nil
	if got := executableName(); got != defaultName {
		t.Errorf("executableName() with no argv = %q, want %q", got, defaultName)
	}
}

// The generated completion script must register the invoked name, since the
// binary is commonly installed as something other than "twitch".
func TestCompletionUsesInvokedName(t *testing.T) {
	orig := rootCmd.Use
	t.Cleanup(func() { rootCmd.Use = orig })

	rootCmd.Use = "ttv"
	out, err := executeCommand("completion", "zsh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "#compdef ttv") {
		t.Errorf("completion should register 'ttv', got first line: %q", firstLine(out))
	}
	if strings.Contains(out, "#compdef twitch") {
		t.Errorf("completion still registers 'twitch': %q", firstLine(out))
	}
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}

// ---- output helpers ----

func TestPrintJSON(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := printJSON(buf, []entities.Category{{ID: "1", Name: "Go"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"id": "1"`, `"name": "Go"`} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("printJSON output missing %q: %s", want, buf.String())
		}
	}
}

func TestPrintStreamTable(t *testing.T) {
	buf := new(bytes.Buffer)
	printStreamTable(buf, []entities.Stream{
		{Username: "streamer1", GameName: "Go", Title: "Coding", ViewerCount: 12},
	})
	for _, want := range []string{"CHANNEL", "streamer1", "Go", "Coding", "12"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("table missing %q:\n%s", want, buf.String())
		}
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		in   string
		max  int
		want string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is far too long", 10, "this is f…"},
		{"héllo wörld", 6, "héllo…"},
	}
	for _, tc := range cases {
		if got := truncate(tc.in, tc.max); got != tc.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tc.in, tc.max, got, tc.want)
		}
	}
}

func TestDashIfEmpty(t *testing.T) {
	if got := dashIfEmpty(""); got != "-" {
		t.Errorf("dashIfEmpty(\"\") = %q, want %q", got, "-")
	}
	if got := dashIfEmpty("x"); got != "x" {
		t.Errorf("dashIfEmpty(\"x\") = %q, want %q", got, "x")
	}
}

// ---- player selection ----

func TestIsBrowserChoice(t *testing.T) {
	for _, in := range []string{"browser", "Browser", "BROWSER"} {
		if !isBrowserChoice(in) {
			t.Errorf("isBrowserChoice(%q) = false, want true", in)
		}
	}
	for _, in := range []string{"", "mpv", "browsers"} {
		if isBrowserChoice(in) {
			t.Errorf("isBrowserChoice(%q) = true, want false", in)
		}
	}
}

func TestMatchesStream(t *testing.T) {
	s := entities.Stream{Username: "GoGopher", GameName: "Software and Game Dev", Title: "Building a CLI"}

	for _, in := range []string{"gopher", "SOFTWARE", "cli"} {
		if !matchesStream(&s, in) {
			t.Errorf("matchesStream(%q) = false, want true", in)
		}
	}
	if matchesStream(&s, "minecraft") {
		t.Error("matchesStream(\"minecraft\") = true, want false")
	}
}

// ---- notifications ----

// Stream titles routinely contain quotes and backslashes; they must not be
// able to terminate or alter the AppleScript literal.
func TestAppleScriptString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`plain`, `"plain"`},
		{`say "hi"`, `"say \"hi\""`},
		{`back\slash`, `"back\\slash"`},
	}
	for _, tc := range cases {
		if got := appleScriptString(tc.in); got != tc.want {
			t.Errorf("appleScriptString(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
