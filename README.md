# gsh

> A Go shell that means business. Sharp prompt, colour themes, and the full weight of POSIX commands in one binary.

![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS-lightgrey)

---

## What it is

**gsh** is a working shell you can compile, install, and use as your daily driver or as a sandboxed sub-shell.

- `ls`, `git`, `python`, `docker`, `curl`: everything in your `$PATH` works
- Pipes (`|`), redirections (`>`, `>>`, `<`), subshells (`$(...)`) all work
- Environment variables (`$HOME`, custom exports) expand correctly
- Compiles to a **single binary**: no Node, no browser, no runtime

---

## Features

| Feature | Details |
|---------|---------|
| **Real command execution** | Full POSIX sh semantics via `mvdan.cc/sh/v3` |
| **Pipes and redirects** | `ls \| grep go > files.txt` works out of the box |
| **Clean prompt** | Username, CWD (auto-truncated), git branch, exit-code indicator |
| **Colour themes** | Dracula, Nord, Catppuccin, One Dark, Tokyo Night |
| **Command history** | Persisted to `~/.config/gsh/history`, arrow-key navigation |
| **Tab completion** | Executables on `$PATH` and file/directory paths |
| **Built-in commands** | `cd`, `export`, `unset`, `env`, `history`, `theme`, `clear`, `help` |
| **TOML config** | Auto-created at `~/.config/gsh/config.toml` on first run |

---

## Installation

### Prerequisites

- Go 1.18 or newer (`go version`)
- Git

### Build from source

```bash
git clone https://github.com/jk08y/gsh.git
cd gsh
go build -o gsh .
```

That produces a `gsh` binary in the current directory.

### Install system-wide (optional)

```bash
sudo mv gsh /usr/local/bin/
```

Then run:

```bash
gsh
```

---

## Usage

```
$ gsh

  gsh  v1.0.0
  theme: dracula  |  type 'help' for built-ins
  All commands run for real. Pipes, redirects, env vars all work.

jk ~/projects/gsh  main
$ _
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
| `version` | Print gsh version |
| `clear` | Clear the screen |
| `exit` / `quit` | Exit the shell |
| `help` | Show built-in help |

### Themes

```bash
theme            # list all available themes
theme nord
theme catppuccin
theme onedark
theme tokyo-night
theme dracula
```

---

## Configuration

On first run gsh creates `~/.config/gsh/config.toml`:

```toml
# gsh configuration

theme     = "dracula"   # dracula | nord | catppuccin | onedark | tokyo-night
show_git  = true        # show git branch in prompt
show_user = true        # show username in prompt
show_host = false       # show hostname in prompt

[history]
max_size = 10000
file     = "~/.config/gsh/history"
```

Edit this file and restart gsh to apply changes.

---

## Project layout

```
gsh/
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
| [`mvdan.cc/sh/v3`](https://github.com/mvdan/sh) | POSIX shell interpreter |
| [`chzyer/readline`](https://github.com/chzyer/readline) | Line editor with history and tab completion |
| [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) | Terminal colour and style rendering |
| [`BurntSushi/toml`](https://github.com/BurntSushi/toml) | Config file parsing |

---

## Contributing

Issues and pull requests are welcome.
