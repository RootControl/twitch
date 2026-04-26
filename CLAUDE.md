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

**Prerequisite:** The [Twitch CLI](https://dev.twitch.tv/docs/cli/) must be installed and authenticated (`twitch token`) for API calls to work. All API calls are executed as shell subprocesses via `twitch api get <endpoint>`.

## Architecture

This is a Go CLI tool that wraps the Twitch CLI to query the Twitch API.

**Request flow:** `main.go` → `api/` functions → `api.Request.Get()` → shells out to `twitch api get <endpoint> <flags>` → JSON unmarshalled into response structs → entity methods format output.

**Key design decision:** Rather than calling the Twitch REST API directly with HTTP, every API call is delegated to the Twitch CLI binary via `exec.Command("sh", "-c", ...)`. The `Request` struct in [api/twitchAPI.go](api/twitchAPI.go) builds the CLI command string and runs it, capturing stdout as a `bytes.Buffer` for JSON unmarshalling.

**Packages:**
- `api/` — one file per resource type (`streamers.go`, `users.go`, `categories.go`), each defining a response struct and one or more query functions. Shared CLI invocation logic lives in `twitchAPI.go`.
- `entities/` — plain data structs with JSON tags mirroring Twitch API shapes, plus formatting methods (`ToString()`, `GetMainInfo()`).
- `browser/` — cross-platform `OpenBrowser(url)` helper (currently unused in `main.go`).

**Adding a new endpoint:** create a file in `api/`, define a response struct embedding the relevant entity types, call `NewApiRequest().Get("<twitch-api-path>", "-q param=value", ...)`, and unmarshal the buffer.

**Known TODO:** `GetFollowedStreams` has a hardcoded `user_id`; the plan is to read it from `.env` (via `godotenv`) or fall back to `GetUser`.
