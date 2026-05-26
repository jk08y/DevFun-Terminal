# gsh

> A shell written in Go. POSIX-compatible, themeable, with aliases, persistent history, and a git-aware prompt.

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
| **Aliases** | Define, persist, and expand shell aliases |
| **RC file** | `~/.config/gsh/gshrc` executed on every startup |
| **Source files** | `source file` runs any script in the current shell context |
| **Directory stack** | `pushd`, `popd`, `dirs` for fast directory navigation |
| **Command lookup** | `which` and `type` to inspect commands and aliases |
| **Clean prompt** | Username, CWD (auto-truncated), git branch, exit-code indicator |
| **Colour themes** | Dracula, Nord, Catppuccin, One Dark, Tokyo Night |
| **Command history** | Persisted to `~/.config/gsh/history`, arrow-key navigation |
| **Tab completion** | Executables on `$PATH` and file/directory paths |
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

---

## Built-in commands

| Command | Description |
|---------|-------------|
| `alias [name[=value]]` | Define or show aliases |
| `unalias name` | Remove an alias |
| `cd [dir]` | Change directory (`cd -` goes back) |
| `pushd [dir]` | Push directory onto stack and cd |
| `popd` | Pop directory stack and cd back |
| `dirs` | Show directory stack |
| `source file` | Execute file in current shell context |
| `. file` | Same as source |
| `which command` | Locate a command or show its alias |
| `type name` | Describe what a name is (builtin, alias, or path) |
| `export KEY=VAL` | Set environment variable |
| `unset KEY` | Unset environment variable |
| `env` | Print all environment variables |
| `history` | Show command history |
| `theme [name]` | List or switch colour theme |
| `version` | Print gsh version |
| `clear` | Clear the screen |
| `exit` / `quit` | Exit the shell |
| `help` | Show built-in help |

---

## Aliases

Aliases are defined with `alias`, persisted automatically to `~/.config/gsh/aliases`, and loaded on every startup.

```bash
alias ll='ls -la'
alias gs='git status'
alias ..='cd ..'

alias          # list all defined aliases
alias ll       # show a single alias
unalias ll     # remove an alias
```

You can also pre-load aliases in your RC file (see below).

---

## RC file

On startup gsh executes `~/.config/gsh/gshrc` if it exists. Use it for aliases, exports, and any setup commands:

```bash
# ~/.config/gsh/gshrc

alias ll='ls -la'
alias gs='git status'
alias gp='git push'

export EDITOR=vim
export GOPATH=$HOME/go
```

---

## Directory stack

```bash
pushd ~/projects   # cd to ~/projects and save current dir
pushd /tmp         # cd to /tmp and save ~/projects
dirs               # show stack: /tmp  ~/projects  (original)
popd               # return to ~/projects
popd               # return to original dir
```

---

## Themes

```bash
theme              # list all available themes
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

alias_file = "~/.config/gsh/aliases"
rc_file    = "~/.config/gsh/gshrc"
```

---

## Config files

| File | Purpose |
|------|---------|
| `~/.config/gsh/config.toml` | Main configuration |
| `~/.config/gsh/gshrc` | Startup script (aliases, exports, etc.) |
| `~/.config/gsh/aliases` | Persisted alias definitions |
| `~/.config/gsh/history` | Command history |

---

## Project layout

```
gsh/
├── main.go
├── go.mod / go.sum
└── internal/
    ├── config/config.go       # TOML config loader
    ├── theme/theme.go         # colour palette definitions
    ├── prompt/prompt.go       # prompt renderer
    ├── completer/completer.go # readline tab-completion
    └── shell/
        ├── shell.go           # REPL loop and startup
        ├── builtins.go        # built-in command handlers
        └── aliases.go         # alias load, save, expand
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
