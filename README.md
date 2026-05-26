# nexterm

> A real, compiled shell written in Go — with a beautiful prompt, themes, and full POSIX command support.

![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-lightgrey)

---

## What it is

**nexterm** is a working shell you can compile, install, and use as your daily driver (or as a sandboxed sub-shell). It is **not** an emulator or a toy — every command you type runs for real on your OS.

- `ls`, `git`, `python`, `docker`, `curl` — everything in your `$PATH` works
- Pipes (`|`), redirections (`>`, `>>`, `<`), subshells (`$(...)`) all work
- Environment variables (`$HOME`, custom exports) expand correctly
- Compiles to a **single static-ish binary** — no Node, no browser, no runtime

---

## Features

| Feature | Details |
|---------|---------|
| **Real command execution** | Full POSIX sh semantics via `mvdan.cc/sh/v3` |
| **Pipes & redirects** | `ls \| grep go > files.txt` works out of the box |
| **Beautiful prompt** | Username · CWD (auto-truncated) · git branch · exit-code arrow |
| **Colour themes** | Dracula · Nord · Catppuccin · One Dark · Tokyo Night |
| **Command history** | Persisted to `~/.config/nexterm/history`, arrow-key navigation |
| **Tab completion** | Executables on `$PATH` + file/directory paths |
| **Built-in commands** | `cd`, `export`, `unset`, `env`, `history`, `theme`, `clear`, `help` |
| **TOML config** | Auto-created at `~/.config/nexterm/config.toml` on first run |

---

## Installation

### Prerequisites

- Go 1.18 or newer (`go version`)
- Git (to clone)

### Build from source

```bash
git clone https://github.com/jk08y/nexterm.git
cd nexterm
go build -o nexterm .
```

That produces a `nexterm` binary in the current directory.

### Install system-wide (optional)

```bash
sudo mv nexterm /usr/local/bin/
```

Then just run:

```bash
nexterm
```

---

## Usage

```
$ nexterm

  nexterm  v1.0.0 — a real shell
  theme: dracula  •  type 'help' for built-ins
  All commands run for real. Pipes, redirects, env vars — everything works.

jk ~/projects/nexterm  main
❯ _
```

### Built-in commands

| Command | Description |
|---------|-------------|
| `cd [dir]` | Change directory (`cd -` goes back) |
| `export KEY=VAL` | Set environment variable |
| `unset KEY` | Unset environment variable |
| `env` | Print all environment variables |
| `history` | Show command history |
| `theme [name]` | List or switch colour theme |
| `version` | Print nexterm version |
| `clear` | Clear the screen |
| `exit` / `quit` | Exit the shell |
| `help` | Show built-in help |

### Themes

```bash
theme            # list all available themes (current one is highlighted)
theme nord       # switch to Nord
theme catppuccin # switch to Catppuccin Mocha
theme onedark    # switch to One Dark
theme tokyo-night
theme dracula    # back to default
```

---

## Configuration

On first run nexterm creates `~/.config/nexterm/config.toml`:

```toml
# nexterm configuration

theme     = "dracula"   # dracula | nord | catppuccin | onedark | tokyo-night
show_git  = true        # show git branch in prompt
show_user = true        # show username in prompt
show_host = false       # show hostname in prompt

[history]
max_size = 10000
file     = "~/.config/nexterm/history"
```

Edit this file and restart nexterm to apply changes.

---

## Project layout

```
nexterm/
├── main.go                    # entry point
├── go.mod / go.sum
└── internal/
    ├── config/config.go       # TOML config loader
    ├── theme/theme.go         # colour palette definitions
    ├── prompt/prompt.go       # prompt renderer
    ├── completer/completer.go # readline tab-completion
    └── shell/
        ├── shell.go           # REPL loop
        └── builtins.go        # built-in commands
```

---

## Tech stack

| Package | Purpose |
|---------|---------|
| [`mvdan.cc/sh/v3`](https://github.com/mvdan/sh) | POSIX shell interpreter (real command execution) |
| [`chzyer/readline`](https://github.com/chzyer/readline) | Line editor with history & tab completion |
| [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) | Terminal colour & style rendering |
| [`BurntSushi/toml`](https://github.com/BurntSushi/toml) | Config file parsing |

---

## Contributing

Issues and pull requests are welcome.

## Connect

𝕏 [@jk08y](https://x.com/jk08y)
