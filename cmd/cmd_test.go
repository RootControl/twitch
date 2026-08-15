package cmd

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/RootControl/twitch/config"
	"github.com/RootControl/twitch/entities"
	"github.com/manifoldco/promptui"
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

func TestDisplayWidth(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"plain", 5},
		{"héllo", 5}, // precomposed accents are single-width
		{"🌸", 2},     // emoji occupy two columns
		{"a🌸b", 4},
		{"日本語", 6}, // CJK is double-width
		{"é", 1},  // combining acute adds no width
		{"👍️", 2},  // variation selector adds no width
	}
	for _, tc := range cases {
		if got := displayWidth(tc.in); got != tc.want {
			t.Errorf("displayWidth(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

// Every line handed to promptui must fit the terminal: its screen buffer
// counts logical lines, so a wrapped line corrupts the redraw.
func TestTruncateNeverExceedsWidth(t *testing.T) {
	inputs := []string{
		"short",
		"🌸🌸🌸🌸🌸🌸🌸🌸🌸🌸🌸🌸",
		"FREAK BOB FRIDAY | #BUNGULATE 🌸 [English funny gaming Ptuber fpsgame]",
		"日本語のタイトルです",
		"mixed 🌸 content with 日本語 and ascii",
	}
	for _, in := range inputs {
		for _, max := range []int{1, 2, 5, 10, 20, 40} {
			got := truncate(in, max)
			if w := displayWidth(got); w > max {
				t.Errorf("truncate(%q, %d) = %q, width %d exceeds %d", in, max, got, w, max)
			}
		}
	}
}

func TestTruncateZeroWidth(t *testing.T) {
	if got := truncate("anything", 0); got != "" {
		t.Errorf("truncate(_, 0) = %q, want empty", got)
	}
}

// The budgets must keep a rendered list row inside the terminal.
func TestItemWidthsFitTerminal(t *testing.T) {
	for _, width := range []int{20, 40, 80, 120, 200} {
		name, game, detail := itemWidths(width)

		// "▶ NAME (GAME) - 12345 viewers"
		row := displayWidth("▶ ") + name + displayWidth(" (") + game +
			displayWidth(") - ") + len("12345") + displayWidth(" viewers")
		if width >= 60 && row > width {
			t.Errorf("width %d: list row needs %d columns", width, row)
		}
		// "Channel:  VALUE"
		if width >= 60 && 10+detail > width {
			t.Errorf("width %d: detail row needs %d columns", width, 10+detail)
		}
		if name <= 0 || game <= 0 || detail <= 0 {
			t.Errorf("width %d: got non-positive budgets (%d, %d, %d)", width, name, game, detail)
		}
	}
}

func TestNewStreamItemsTruncates(t *testing.T) {
	long := strings.Repeat("very long title ", 20)
	items := newStreamItems([]entities.Stream{
		{Username: "streamer1", GameName: "Software and Game Dev", Title: long, ViewerCount: 5,
			StartedAt: "2020-01-01T00:00:00Z"},
	}, 80)

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	it := items[0]
	if displayWidth(it.Title) > 80 {
		t.Errorf("title not truncated: width %d", displayWidth(it.Title))
	}
	if it.Viewers != 5 {
		t.Errorf("Viewers = %d, want 5", it.Viewers)
	}
	if it.Uptime == "" {
		t.Error("Uptime should be populated for a stream with a start time")
	}
	// The untruncated stream must be preserved for searching and opening.
	if it.stream.Title != long {
		t.Error("original stream title should be preserved on the item")
	}
}

var ansi = regexp.MustCompile("\x1b\\[[0-9;]*m")

// This is the end-to-end guard for the redraw corruption: render the real
// templates against a hostile stream and assert every line fits the terminal.
// A single over-wide line makes promptui wrap, miscount its own height, and
// leave the previous frame on screen when the user moves the cursor.
func TestRenderedLinesFitTerminal(t *testing.T) {
	stream := entities.Stream{
		Username:    "AVeryLongChannelNameThatKeepsGoingAndGoing",
		GameName:    "A Category With An Unreasonably Long Name Indeed",
		Title:       "First time playing Grand Theft Auto V!!! (Day 3) 🌸 | @marzzzzy !socials !fractalSS | FREAK BOB FRIDAY | #BUNGULATE [English funny gaming Ptuber fpsgame]",
		ViewerCount: 123456,
		StartedAt:   "2020-01-01T00:00:00Z",
	}

	tmpls := streamSelectTemplates()
	for _, width := range []int{60, 80, 100, 120, 200} {
		items := newStreamItems([]entities.Stream{stream}, width)

		for name, src := range map[string]string{
			"active":   tmpls.Active,
			"inactive": tmpls.Inactive,
			"details":  tmpls.Details,
		} {
			tmpl, err := template.New(name).Funcs(promptui.FuncMap).Parse(src)
			if err != nil {
				t.Fatalf("parsing %s: %v", name, err)
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, items[0]); err != nil {
				t.Fatalf("executing %s: %v", name, err)
			}

			for _, line := range strings.Split(buf.String(), "\n") {
				// Colour codes occupy no columns on screen.
				plain := ansi.ReplaceAllString(line, "")
				if w := displayWidth(plain); w > width {
					t.Errorf("width %d: %s line is %d columns: %q", width, name, w, plain)
				}
			}
		}
	}
}

// The details pane must show a real uptime value, not a dump of the struct.
// The underlying pointer-receiver bug is covered by
// entities.TestStreamMethodsResolveFromTemplate; this asserts the pane the
// user actually sees.
func TestDetailsRendersUptime(t *testing.T) {
	items := newStreamItems([]entities.Stream{
		{Username: "streamer1", StartedAt: "2020-01-01T00:00:00Z"},
	}, 80)

	tmpl, err := template.New("d").Funcs(promptui.FuncMap).Parse(streamSelectTemplates().Details)
	if err != nil {
		t.Fatalf("parsing details: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, items[0]); err != nil {
		t.Fatalf("executing details: %v", err)
	}

	out := ansi.ReplaceAllString(buf.String(), "")
	if !strings.Contains(out, "Uptime:") {
		t.Fatalf("details missing the uptime label:\n%s", out)
	}
	// A struct dump would carry the thumbnail URL and the language field.
	if strings.Contains(out, "https://static-cdn") || strings.Contains(out, "{") {
		t.Errorf("details rendered a raw struct:\n%s", out)
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
