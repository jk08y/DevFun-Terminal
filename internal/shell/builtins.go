package shell

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jk08y/nexterm/internal/theme"
)

// builtinResult is returned by every builtin handler.
type builtinResult struct {
	handled  bool // was this command a builtin?
	exitCode int
	doExit   bool // should the shell terminate?
}

// handleBuiltin checks if args[0] is a builtin and runs it.
// Returns (result, true) when the command was consumed.
func (s *Shell) handleBuiltin(args []string) (builtinResult, bool) {
	if len(args) == 0 {
		return builtinResult{}, false
	}

	switch args[0] {
	case "cd":
		return s.builtinCd(args[1:]), true
	case "exit", "quit":
		return builtinResult{handled: true, doExit: true}, true
	case "clear":
		fmt.Print("\033[H\033[2J")
		return builtinResult{handled: true}, true
	case "history":
		return s.builtinHistory(), true
	case "theme":
		return s.builtinTheme(args[1:]), true
	case "help":
		return s.builtinHelp(), true
	case "export":
		return s.builtinExport(args[1:]), true
	case "unset":
		return s.builtinUnset(args[1:]), true
	case "env":
		return s.builtinEnv(), true
	case "version":
		return s.builtinVersion(), true
	}

	return builtinResult{}, false
}

// ── cd ──────────────────────────────────────────────────────────────────────

func (s *Shell) builtinCd(args []string) builtinResult {
	var dir string
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: cannot find home directory\n")
			return builtinResult{handled: true, exitCode: 1}
		}
		dir = home
	} else if args[0] == "-" {
		dir = os.Getenv("OLDPWD")
		if dir == "" {
			fmt.Fprintf(os.Stderr, "cd: OLDPWD not set\n")
			return builtinResult{handled: true, exitCode: 1}
		}
		fmt.Println(dir)
	} else {
		dir = args[0]
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		return builtinResult{handled: true, exitCode: 1}
	}
	_ = os.Setenv("OLDPWD", prev)
	newDir, _ := os.Getwd()
	_ = os.Setenv("PWD", newDir)

	return builtinResult{handled: true}
}

// ── history ──────────────────────────────────────────────────────────────────

func (s *Shell) builtinHistory() builtinResult {
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)
	for i, cmd := range s.history {
		fmt.Printf("%s  %s\n", muted.Render(fmt.Sprintf("%4d", i+1)), cmd)
	}
	return builtinResult{handled: true}
}

// ── theme ─────────────────────────────────────────────────────────────────────

func (s *Shell) builtinTheme(args []string) builtinResult {
	if len(args) == 0 {
		// List available themes
		available := theme.List()
		sort.Strings(available)
		primaryStyle := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
		fmt.Println("Available themes:")
		for _, name := range available {
			marker := "  "
			if name == s.cfg.Theme {
				marker = primaryStyle.Render("▶ ")
			}
			fmt.Println(marker + name)
		}
		return builtinResult{handled: true}
	}

	name := args[0]
	t := theme.Get(name)
	if t.Name != name {
		fmt.Fprintf(os.Stderr, "theme: unknown theme %q — run 'theme' to list available\n", name)
		return builtinResult{handled: true, exitCode: 1}
	}

	s.theme = t
	s.cfg.Theme = name
	s.promptBuilder.SetTheme(t)

	ok := lipgloss.NewStyle().Foreground(t.Success).Bold(true)
	fmt.Printf("%s switched to theme %q\n", ok.Render("✓"), name)
	return builtinResult{handled: true}
}

// ── export ───────────────────────────────────────────────────────────────────

func (s *Shell) builtinExport(args []string) builtinResult {
	for _, arg := range args {
		k, v, found := strings.Cut(arg, "=")
		if !found {
			// export NAME with existing value — just mark as exported (already is on Linux)
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			fmt.Fprintf(os.Stderr, "export: %v\n", err)
			return builtinResult{handled: true, exitCode: 1}
		}
	}
	return builtinResult{handled: true}
}

// ── unset ─────────────────────────────────────────────────────────────────────

func (s *Shell) builtinUnset(args []string) builtinResult {
	for _, key := range args {
		_ = os.Unsetenv(key)
	}
	return builtinResult{handled: true}
}

// ── env ──────────────────────────────────────────────────────────────────────

func (s *Shell) builtinEnv() builtinResult {
	for _, kv := range os.Environ() {
		fmt.Println(kv)
	}
	return builtinResult{handled: true}
}

// ── version ──────────────────────────────────────────────────────────────────

func (s *Shell) builtinVersion() builtinResult {
	bold := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
	fmt.Printf("%s v%s\n", bold.Render("nexterm"), s.version)
	return builtinResult{handled: true}
}

// ── help ──────────────────────────────────────────────────────────────────────

func (s *Shell) builtinHelp() builtinResult {
	title := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
	cmd := lipgloss.NewStyle().Foreground(s.theme.Secondary)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)

	fmt.Printf("\n%s — built-in commands\n\n", title.Render("nexterm"))

	builtins := []struct{ name, desc string }{
		{"cd [dir]", "change directory (cd - goes back)"},
		{"export KEY=VAL", "set environment variable"},
		{"unset KEY", "unset environment variable"},
		{"env", "print all environment variables"},
		{"history", "show command history"},
		{"theme [name]", "list or switch colour theme"},
		{"version", "print nexterm version"},
		{"clear", "clear the screen"},
		{"exit / quit", "exit the shell"},
		{"help", "show this message"},
	}

	for _, b := range builtins {
		fmt.Printf("  %-28s %s\n", cmd.Render(b.name), muted.Render(b.desc))
	}

	fmt.Println()
	fmt.Println(muted.Render("All other input is executed as a real shell command (pipes, redirects, etc. work)."))
	fmt.Println()

	return builtinResult{handled: true}
}
