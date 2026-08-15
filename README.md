# twitch

A Go CLI for browsing Twitch from the terminal — list live streams, search channels, follow broadcasters, and open a stream in `streamlink`, `mpv`, or your browser.

It is a wrapper around the [official Twitch CLI](https://dev.twitch.tv/docs/cli/): every request is delegated to the `twitch` binary, so authentication is handled once by `twitch token` and never stored by this tool.

## Requirements

- Go 1.24 or newer
- The [Twitch CLI](https://dev.twitch.tv/docs/cli/), installed and authenticated
- Optional: [`streamlink`](https://streamlink.github.io/) or [`mpv`](https://mpv.io/) to play streams outside the browser

```bash
brew install twitchdev/twitch/twitch-cli   # or see the docs for your platform
twitch configure                           # enter your client ID and secret
twitch token                               # app token, for browsing
```

The `follow`, `unfollow`, `followed`, and `watch` commands act on your account, so they need a **user token** with the right scopes:

```bash
twitch token -u -s 'user:read:follows user:edit:follows'
```

## Install

```bash
git clone https://github.com/RootControl/twitch.git
cd twitch
make install
```

That builds the binary as `ttv` and installs it to `$GOBIN` (or `$(go env GOPATH)/bin`), which the Go toolchain already expects on your `PATH`. `make install` warns if it isn't. For a system-wide install:

```bash
sudo make install PREFIX=/usr/local/bin
```

> **Do not install this binary as `twitch`.** `go build` names it `twitch` by default, which collides with the Twitch CLI it shells out to. If it lands earlier on your `PATH` than the real CLI, it will invoke itself recursively. Override the name only if you have somewhere safe for it: `make install BINARY=tw`.

The tool adopts whatever name you invoke it as, so its help text and completion scripts follow suit. The examples below use `twitch` for readability — substitute the name you installed.

Run `make help` for all targets. The useful ones:

| Target | Description |
| --- | --- |
| `make` | Build `./ttv` |
| `make install` | Build and install to `$(PREFIX)` |
| `make uninstall` | Remove the installed binary |
| `make run ARGS="streams -g Rust"` | Build and run |
| `make check` | gofmt check, `go vet`, and tests — run before pushing |
| `make test` / `make race` | Tests, with or without the race detector |
| `make clean` | Remove build artifacts |

## Configuration

No configuration is required. Commands that act on your account resolve your user ID automatically from the token the Twitch CLI is authenticated with.

You can skip that lookup, or point the tool at a different account, with a `.env` file (see [.env.example](.env.example)):

```bash
TWITCH_USER_ID=123456789      # optional; skips one API call
TWITCH_USER_LOGIN=your_login  # optional; resolved to an ID if USER_ID is unset
```

## Commands

### Browsing streams

#### `streams` — list and open top live streams

```bash
twitch streams                          # top 30 live streams overall
twitch streams -g "just chatting"       # filter by category
twitch streams -g Rust -n 50            # 50 results
twitch streams --json                   # machine-readable
```

| Flag | Description |
| --- | --- |
| `-g, --game` | Filter by game/category name. Resolved to a category ID; an exact name match wins over a fuzzy one. |
| `-n, --limit` | Number of streams to fetch, 1–100 (default 30) |
| `-p, --player` | Open directly with `streamlink`, `mpv`, or `browser`, skipping the prompt |
| `--json` | Output as JSON |

On a terminal this opens an interactive picker — arrow keys to move, `/` to filter by channel, category, or title, `Enter` to select. After picking a stream you choose a player. When output is piped or redirected the picker is skipped and a plain table is printed instead:

```
$ twitch streams -n 3
CHANNEL         VIEWERS  UPTIME  CATEGORY            TITLE
TheBurntPeanut  48371    2h57m   Escape from Tarkov  [DROPS] TARKOV SEASONAL | TEACHING TIMTHET…
dona            26485    3h06m   EA Sports FC 26     AMISTOSO - ⚽ | ICE NUGGETS (0) x (0) LOLO…
AussieAntics    24356    8h45m   Fortnite            DROPS ON 🔴 WATCHING FNCS LAST CHANCE FINA…
```

#### `followed` — list and open followed live streams

```bash
twitch followed
twitch followed --json          # prints [] when nothing is live
twitch followed -p mpv          # open the pick in mpv
```

Same flags as `streams`, minus `--game`. Requires a user token.

#### `open` — open a channel directly

```bash
twitch open shroud
twitch open shroud -p streamlink
```

Verifies the channel exists before launching anything, so a typo gives you a clear message rather than a 404 page.

### Discovery

#### `search` — find channels by name

Matches channel names and returns both live and offline broadcasters.

```bash
twitch search shroud
twitch search shroud --live      # currently-live only
twitch search shroud -n 5 --json
```

```
$ twitch search shroud --live -n 1
LIVE     shrood
         24/7 @SHROUD | VODS AND STREAMS | CS:GO | PUBG | VALO | RPGS | REPLAYS
         Always On
         https://www.twitch.tv/shrood
```

| Flag | Description |
| --- | --- |
| `-l, --live` | Only channels that are currently live |
| `-n, --limit` | Number of channels to fetch, 1–100 (default 20) |
| `--json` | Output as JSON |

#### `top` — most-watched categories right now

```bash
twitch top
twitch top -n 5
```

```
$ twitch top -n 5
 1. Just Chatting
 2. Escape from Tarkov
 3. Fortnite
 4. Grand Theft Auto V
 5. IRL
```

#### `categories` — search categories and games

Prints matching categories with their IDs.

```bash
twitch categories "just chatting"
twitch categories rust -n 10 --json
```

```
$ twitch categories "just chatting" -n 1
ID: 509658	Name: Just Chatting
```

### Users and channels

#### `user` — user profile

```bash
twitch user shroud
twitch user shroud --json
```

#### `channel` — channel info

Shows the current title, category, language, and link.

```bash
twitch channel shroud
```

```
Channel:  shroud
Title:    ME N THE GIRLS MAY BE OUT! BUT WE HAD FUN!
Category: VALORANT
Language: en
Link:     https://www.twitch.tv/shroud
```

### Following

Both require a user token with `user:edit:follows`.

```bash
twitch follow shroud
twitch unfollow shroud
```

### `watch` — notify when a followed channel goes live

Polls your followed streams and prints a line (plus a desktop notification) when someone starts broadcasting. Channels already live when `watch` starts are treated as the baseline, not as new.

```bash
twitch watch                    # poll every 60s
twitch watch -i 300             # every 5 minutes
twitch watch --notify=false     # terminal output only
```

| Flag | Description |
| --- | --- |
| `-i, --interval` | Polling interval in seconds, minimum 10 (default 60) |
| `--notify` | Send a desktop notification (default true) |

```
$ twitch watch -i 120
Watching followed streams (3 live now, polling every 120s). Press Ctrl+C to stop.
[21:14:07] Marzzzzy just went live: First time playing GTA V!!! (Grand Theft Auto V)
```

Desktop notifications use `osascript` on macOS and `notify-send` on Linux; Windows prints to the terminal only. A failed poll is reported and retried on the next tick rather than ending the session. `Ctrl+C` exits cleanly.

### `completion` — shell autocompletion

Supports `bash`, `zsh`, `fish`, and `powershell`.

```bash
ttv completion zsh > "${fpath[1]}/_ttv"
ttv completion bash > /etc/bash_completion.d/ttv
```

The script registers whatever name you invoked the binary as, so it matches the name you installed it under with no editing.

## Scripting

Every listing command takes `--json` and writes to stdout, so results compose with other tools:

```bash
twitch followed --json | jq -r '.[] | "\(.user_name)\t\(.viewer_count)"'
twitch top --json | jq -r '.[0].name'
twitch streams -g "Software and Game Dev" --json | jq 'length'
```

Errors go to stderr and set a non-zero exit status, and API failures report the message Twitch returned:

```
$ twitch user no-such-user-abc-123
Error: user "no-such-user-abc-123" not found

$ twitch followed          # with an expired token
Error: twitch API 401: Invalid OAuth token (run `twitch token` to re-authenticate)
```

## Development

```bash
go build ./...    # build all packages
go vet ./...      # static analysis
go test ./...     # run all tests
go test ./api/... # one package
```

No test hits the network: each `api` function has an unexported twin that takes an `Executor`, so tests inject a fake that captures the arguments and returns a canned payload.

See [CLAUDE.md](CLAUDE.md) for the architecture and for conventions to follow when adding an endpoint.
