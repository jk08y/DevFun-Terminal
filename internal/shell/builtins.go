package shell

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jk08y/gsh/internal/theme"
)

// builtinResult is returned by every builtin handler.
type builtinResult struct {
	handled  bool
	exitCode int
	doExit   bool
}

// builtinNames is the authoritative list of gsh built-in command names.
var builtinNames = []string{
	"alias", "unalias", "cd", "clear", "dirs", "env", "env",
	"exec", "exit", "export", "help", "history", "popd", "pushd",
	"quit", "source", "theme", "type", "unset", "version", "which", ".",
}

// isBuiltin reports whether name is a gsh built-in.
func isBuiltin(name string) bool {
	for _, b := range builtinNames {
		if b == name {
			return true
		}
	}
	return false
}

// handleBuiltin dispatches to the correct handler.
// Returns (result, true) when the command was consumed.
func (s *Shell) handleBuiltin(args []string) (builtinResult, bool) {
	if len(args) == 0 {
		return builtinResult{}, false
	}
	switch args[0] {
	case "alias":
		return s.builtinAlias(args[1:]), true
	case "unalias":
		return s.builtinUnalias(args[1:]), true
	case "cd":
		return s.builtinCd(args[1:]), true
	case "clear":
		fmt.Print("\033[H\033[2J")
		return builtinResult{handled: true}, true
	case "dirs":
		return s.builtinDirs(), true
	case "env":
		return s.builtinEnv(), true
	case "exit", "quit":
		return builtinResult{handled: true, doExit: true}, true
	case "export":
		return s.builtinExport(args[1:]), true
	case "help":
		return s.builtinHelp(), true
	case "history":
		return s.builtinHistory(), true
	case "popd":
		return s.builtinPopd(), true
	case "pushd":
		return s.builtinPushd(args[1:]), true
	case "source", ".":
		return s.builtinSource(args[1:]), true
	case "theme":
		return s.builtinTheme(args[1:]), true
	case "type":
		return s.builtinType(args[1:]), true
	case "unset":
		return s.builtinUnset(args[1:]), true
	case "version":
		return s.builtinVersion(), true
	case "which":
		return s.builtinWhich(args[1:]), true
	}
	return builtinResult{}, false
}

// ── alias ────────────────────────────────────────────────────────────────────

func (s *Shell) builtinAlias(args []string) builtinResult {
	secondary := lipgloss.NewStyle().Foreground(s.theme.Secondary)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)

	if len(args) == 0 {
		if len(s.aliases) == 0 {
			fmt.Println(muted.Render("no aliases defined"))
			return builtinResult{handled: true}
		}
		keys := make([]string, 0, len(s.aliases))
		for k := range s.aliases {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("%s '%s'\n", secondary.Render("alias "+k+"="), s.aliases[k])
		}
		return builtinResult{handled: true}
	}

	for _, arg := range args {
		name, value, found := strings.Cut(arg, "=")
		if !found {
			// Show single alias
			if v, ok := s.aliases[name]; ok {
				fmt.Printf("%s '%s'\n", secondary.Render("alias "+name+"="), v)
			} else {
				fmt.Fprintf(os.Stderr, "alias: %s: not found\n", name)
				return builtinResult{handled: true, exitCode: 1}
			}
			continue
		}
		s.aliases[name] = value
	}
	s.saveAliases()
	return builtinResult{handled: true}
}

// ── unalias ──────────────────────────────────────────────────────────────────

func (s *Shell) builtinUnalias(args []string) builtinResult {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "unalias: usage: unalias name [name...]\n")
		return builtinResult{handled: true, exitCode: 1}
	}
	for _, name := range args {
		delete(s.aliases, name)
	}
	s.saveAliases()
	return builtinResult{handled: true}
}

// ── cd ───────────────────────────────────────────────────────────────────────

func (s *Shell) builtinCd(args []string) builtinResult {
	var dir string
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "cd: cannot find home directory")
			return builtinResult{handled: true, exitCode: 1}
		}
		dir = home
	} else if args[0] == "-" {
		dir = os.Getenv("OLDPWD")
		if dir == "" {
			fmt.Fprintln(os.Stderr, "cd: OLDPWD not set")
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

// ── pushd / popd / dirs ──────────────────────────────────────────────────────

func (s *Shell) builtinPushd(args []string) builtinResult {
	cur, _ := os.Getwd()

	if len(args) == 0 {
		// Swap top two entries (like bash)
		if len(s.dirStack) == 0 {
			fmt.Fprintln(os.Stderr, "pushd: directory stack empty")
			return builtinResult{handled: true, exitCode: 1}
		}
		top := s.dirStack[len(s.dirStack)-1]
		s.dirStack[len(s.dirStack)-1] = cur
		if err := os.Chdir(top); err != nil {
			fmt.Fprintf(os.Stderr, "pushd: %v\n", err)
			return builtinResult{handled: true, exitCode: 1}
		}
		newDir, _ := os.Getwd()
		_ = os.Setenv("PWD", newDir)
		return s.builtinDirs()
	}

	dir := args[0]
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "pushd: %v\n", err)
		return builtinResult{handled: true, exitCode: 1}
	}
	s.dirStack = append(s.dirStack, cur)
	newDir, _ := os.Getwd()
	_ = os.Setenv("PWD", newDir)
	return s.builtinDirs()
}

func (s *Shell) builtinPopd() builtinResult {
	if len(s.dirStack) == 0 {
		fmt.Fprintln(os.Stderr, "popd: directory stack empty")
		return builtinResult{handled: true, exitCode: 1}
	}
	top := s.dirStack[len(s.dirStack)-1]
	s.dirStack = s.dirStack[:len(s.dirStack)-1]
	if err := os.Chdir(top); err != nil {
		fmt.Fprintf(os.Stderr, "popd: %v\n", err)
		return builtinResult{handled: true, exitCode: 1}
	}
	newDir, _ := os.Getwd()
	_ = os.Setenv("PWD", newDir)
	return s.builtinDirs()
}

func (s *Shell) builtinDirs() builtinResult {
	cwd, _ := os.Getwd()
	secondary := lipgloss.NewStyle().Foreground(s.theme.Secondary)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)
	fmt.Printf("%s", secondary.Render(cwd))
	for i := len(s.dirStack) - 1; i >= 0; i-- {
		fmt.Printf("  %s", muted.Render(s.dirStack[i]))
	}
	fmt.Println()
	return builtinResult{handled: true}
}

// ── source / . ───────────────────────────────────────────────────────────────

func (s *Shell) builtinSource(args []string) builtinResult {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "source: usage: source file")
		return builtinResult{handled: true, exitCode: 1}
	}
	code := s.sourceFile(args[0], false)
	return builtinResult{handled: true, exitCode: code}
}

// ── which ────────────────────────────────────────────────────────────────────

func (s *Shell) builtinWhich(args []string) builtinResult {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "which: usage: which command [command...]")
		return builtinResult{handled: true, exitCode: 1}
	}
	secondary := lipgloss.NewStyle().Foreground(s.theme.Secondary)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)
	exitCode := 0
	for _, name := range args {
		if isBuiltin(name) {
			fmt.Printf("%s %s\n", secondary.Render(name), muted.Render("(gsh built-in)"))
			continue
		}
		if v, ok := s.aliases[name]; ok {
			fmt.Printf("%s %s\n", secondary.Render(name+":"), muted.Render("aliased to '"+v+"'"))
			continue
		}
		path, err := exec.LookPath(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "which: %s: not found\n", name)
			exitCode = 1
			continue
		}
		fmt.Println(path)
	}
	return builtinResult{handled: true, exitCode: exitCode}
}

// ── type ─────────────────────────────────────────────────────────────────────

func (s *Shell) builtinType(args []string) builtinResult {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "type: usage: type name [name...]")
		return builtinResult{handled: true, exitCode: 1}
	}
	secondary := lipgloss.NewStyle().Foreground(s.theme.Secondary)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)
	exitCode := 0
	for _, name := range args {
		if isBuiltin(name) {
			fmt.Printf("%s is a %s\n", secondary.Render(name), muted.Render("gsh built-in"))
			continue
		}
		if v, ok := s.aliases[name]; ok {
			fmt.Printf("%s is an %s for '%s'\n", secondary.Render(name), muted.Render("alias"), v)
			continue
		}
		path, err := exec.LookPath(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "type: %s: not found\n", name)
			exitCode = 1
			continue
		}
		fmt.Printf("%s is %s\n", secondary.Render(name), path)
	}
	return builtinResult{handled: true, exitCode: exitCode}
}

// ── export ───────────────────────────────────────────────────────────────────

func (s *Shell) builtinExport(args []string) builtinResult {
	for _, arg := range args {
		k, v, found := strings.Cut(arg, "=")
		if !found {
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
		available := theme.List()
		sort.Strings(available)
		primaryStyle := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
		fmt.Println("Available themes:")
		for _, name := range available {
			marker := "  "
			if name == s.cfg.Theme {
				marker = primaryStyle.Render("> ")
			}
			fmt.Println(marker + name)
		}
		return builtinResult{handled: true}
	}
	name := args[0]
	t := theme.Get(name)
	if t.Name != name {
		fmt.Fprintf(os.Stderr, "theme: unknown theme %q  run 'theme' to list available\n", name)
		return builtinResult{handled: true, exitCode: 1}
	}
	s.theme = t
	s.cfg.Theme = name
	s.promptBuilder.SetTheme(t)
	ok := lipgloss.NewStyle().Foreground(t.Success).Bold(true)
	fmt.Printf("%s switched to theme %q\n", ok.Render("ok"), name)
	return builtinResult{handled: true}
}

// ── version ──────────────────────────────────────────────────────────────────

func (s *Shell) builtinVersion() builtinResult {
	bold := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
	fmt.Printf("%s v%s\n", bold.Render("gsh"), s.version)
	return builtinResult{handled: true}
}

// ── help ──────────────────────────────────────────────────────────────────────

func (s *Shell) builtinHelp() builtinResult {
	title := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
	cmd := lipgloss.NewStyle().Foreground(s.theme.Secondary)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)

	fmt.Printf("\n%s built-in commands\n\n", title.Render("gsh"))

	rows := []struct{ name, desc string }{
		{"alias [name[=value]]", "define or show aliases"},
		{"unalias name", "remove an alias"},
		{"cd [dir]", "change directory (cd - goes back)"},
		{"pushd [dir]", "push directory onto stack and cd"},
		{"popd", "pop directory stack and cd back"},
		{"dirs", "show directory stack"},
		{"source file", "execute file in current shell context"},
		{". file", "same as source"},
		{"which command", "locate a command or show its alias"},
		{"type name", "describe what a name is"},
		{"export KEY=VAL", "set environment variable"},
		{"unset KEY", "unset environment variable"},
		{"env", "print all environment variables"},
		{"history", "show command history"},
		{"theme [name]", "list or switch colour theme"},
		{"version", "print gsh version"},
		{"clear", "clear the screen"},
		{"exit / quit", "exit the shell"},
		{"help", "show this message"},
	}

	for _, r := range rows {
		fmt.Printf("  %-30s %s\n", cmd.Render(r.name), muted.Render(r.desc))
	}

	fmt.Println()
	fmt.Println(muted.Render("Everything else is executed as a real OS command. Pipes, redirects, and subshells all work."))
	fmt.Println()
	return builtinResult{handled: true}
}
