# pguard

`pguard` is a simple CLI tool for running commands with optional timeouts and graceful shutdown, similar to Unix `timeout`.

## Installation

```sh
go install github.com/iamBelugax/pguard@latest
````

Or build manually:

```sh
go build -o pguard
```

## Usage

```sh
pguard [flags] <command> [args...]
```

### Flags

* `--timeout` – Maximum runtime (e.g. `10s`, `1m`). `-1` means no timeout.
* `--graceful` – Grace period before force kill (default: `5s`).

## Examples

Run a command with a timeout:

```sh
pguard --timeout=10s sleep 30
```

Gracefully stop a long-running process:

```sh
pguard --timeout=1m --graceful=10s my-command --foo bar
```

## Behavior

1. Starts the specified command.
2. Waits for completion, timeout, or interrupt (`Ctrl+C`).
3. Sends `SIGINT` for graceful shutdown.
4. Force kills the process after the grace period if still running.