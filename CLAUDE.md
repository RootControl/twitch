# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run main.go       # Run the application
go build ./...       # Build all packages
go vet ./...         # Lint/static analysis
go test ./...        # Run all tests
go test ./api/...    # Run tests in a specific package
```

**Prerequisite:** The [Twitch CLI](https://dev.twitch.tv/docs/cli/) must be installed and authenticated (`twitch token`) for API calls to work. All API calls are executed as subprocesses of the `twitch` binary.

## Architecture

This is a Go CLI tool that wraps the Twitch CLI to query the Twitch API.

**Request flow:** `main.go` → `cmd/` (cobra commands) → `api/` functions → `api.fetch()` → runs `twitch api get <endpoint> -q k=v` → JSON unmarshalled into response structs → `entities/` methods format output.

**Key design decision:** Rather than calling the Twitch REST API directly with HTTP, every API call is delegated to the Twitch CLI binary. The `Request` struct in [api/twitchAPI.go](api/twitchAPI.go) builds an **argv slice** and runs `exec.Command("twitch", args...)` — never a shell. Query values are passed as discrete `-q key=value` arguments, so titles and category names containing spaces or shell metacharacters are handed to the CLI verbatim. `Request.ToString()` exists for display and tests only; never feed it to a shell.

**Error handling:** the `api` and `config` packages return errors — they must never call `log.Fatal`, since `watch` is a long-running loop that has to survive a transient failure. Commands use cobra's `RunE`, and `cmd.Execute()` prints the error and exits 1. The Twitch CLI reports API errors as a JSON envelope on **stdout** with a non-zero exit status, so `fetch`/`mutate` check the payload with `checkAPIError` *before* the exit status; otherwise a 401 is indistinguishable from an empty result set.

**Command name:** the default build name (`twitch`) collides with the Twitch CLI this tool shells out to, so it is usually installed under another name. `Execute()` sets `rootCmd.Use` from `os.Args[0]` so usage text and generated completion scripts follow the invoked name. Tests drive `rootCmd.Execute()` directly and therefore keep the `defaultName` literal.

**Packages:**
- `cmd/` — one file per command group, wired up in `root.go`. `output.go` holds shared JSON/table printing and the `isInteractive()` TTY check.
- `api/` — one file per resource type (`streamers.go`, `users.go`, `categories.go`, `channel.go`, `follows.go`), each defining a response struct and query functions. Shared invocation logic lives in `twitchAPI.go`.
- `entities/` — plain data structs with JSON tags mirroring Twitch API shapes, plus formatting methods (`ToString()`, `GetMainInfo()`, `Uptime()`).
- `config/` — reads `TWITCH_USER_ID` / `TWITCH_USER_LOGIN` from `.env` via `godotenv`.
- `player/` — detects installed stream players (streamlink, mpv) and launches them.
- `browser/` — cross-platform `OpenBrowser(url)` helper.

**Testing:** every `api` function has an unexported twin taking an `Executor`, so tests inject a `fakeExecutor` that captures argv and returns a canned payload. No test hits the network. `cmd` tests reuse the `rootCmd` singleton, so `executeCommand` resets flag state first (a leaked `--help` flag will otherwise make later commands no-op).

**User ID resolution:** `api.resolveUserID` resolves the current user in order — `TWITCH_USER_ID`, then a lookup of `TWITCH_USER_LOGIN`, then the user that owns the CLI's token (the `users` endpoint with no params). The result is memoised. `.env` is therefore optional; setting `TWITCH_USER_ID` just skips a lookup.

**Adding a new endpoint:** create a file in `api/`, define a response struct embedding the relevant entity types, and call `fetch[YourResponse](e, "<twitch-api-path>", Q("param", value))`. Take an `Executor` as the last parameter of the unexported form so it can be tested.
