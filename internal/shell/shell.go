package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
	"github.com/jk08y/gsh/internal/completer"
	"github.com/jk08y/gsh/internal/config"
	"github.com/jk08y/gsh/internal/prompt"
	"github.com/jk08y/gsh/internal/theme"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Shell is the central gsh REPL.
type Shell struct {
	cfg           *config.Config
	version       string
	theme         theme.Theme
	promptBuilder *prompt.Builder
	runner        *interp.Runner
	rl            *readline.Instance
	history       []string
	exitCode      int
}

// New creates a Shell from the provided configuration.
func New(cfg *config.Config, version string) *Shell {
	return &Shell{
		cfg:           cfg,
		version:       version,
		theme:         theme.Get(cfg.Theme),
		promptBuilder: prompt.New(cfg),
	}
}

// Cleanup is called on SIGTERM: gracefully closes readline.
func (s *Shell) Cleanup() {
	if s.rl != nil {
		_ = s.rl.Close()
	}
}

// Run starts the interactive REPL loop.
func (s *Shell) Run() error {
	// Ensure history directory exists
	if dir := parentDir(s.cfg.History.File); dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

	// Set up readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 s.buildPrompt(),
		HistoryFile:            s.cfg.History.File,
		HistoryLimit:           s.cfg.History.MaxSize,
		AutoComplete:           completer.New(),
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		return fmt.Errorf("initialising readline: %w", err)
	}
	defer rl.Close()
	s.rl = rl

	// Set up the mvdan.cc/sh runner (real POSIX-compatible execution)
	runner, err := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.Env(nil), // inherit process environment
	)
	if err != nil {
		return fmt.Errorf("initialising shell runner: %w", err)
	}
	s.runner = runner

	s.printWelcome()

	// Main REPL loop
	for {
		rl.SetPrompt(s.buildPrompt())

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			fmt.Println()
			s.exitCode = 130
			continue
		}
		if err == io.EOF {
			fmt.Println()
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Track history in memory (readline persists to file automatically)
		s.history = append(s.history, line)

		// Try builtin first
		args := strings.Fields(os.ExpandEnv(line))
		if result, ok := s.handleBuiltin(args); ok {
			s.exitCode = result.exitCode
			if result.doExit {
				break
			}
			continue
		}

		// External command via the POSIX shell interpreter
		s.exitCode = s.execute(line)
	}

	return nil
}

// execute runs a shell command line (supports pipes, redirects, subshells, etc.)
// and returns its exit code.
func (s *Shell) execute(line string) int {
	f, err := syntax.NewParser().Parse(strings.NewReader(line), "gsh")
	if err != nil {
		errStyle := lipgloss.NewStyle().Foreground(s.theme.Error)
		fmt.Fprintln(os.Stderr, errStyle.Render("parse error: "+err.Error()))
		return 1
	}

	ctx := context.Background()
	if err := s.runner.Run(ctx, f); err != nil {
		if code, ok := interp.IsExitStatus(err); ok {
			return int(code)
		}
		return 1
	}
	return 0
}

// buildPrompt renders the current prompt string.
func (s *Shell) buildPrompt() string {
	return s.promptBuilder.Build(s.exitCode)
}

// printWelcome prints the startup banner.
func (s *Shell) printWelcome() {
	primary := lipgloss.NewStyle().Foreground(s.theme.Primary).Bold(true)
	muted := lipgloss.NewStyle().Foreground(s.theme.Muted)
	secondary := lipgloss.NewStyle().Foreground(s.theme.Secondary)

	fmt.Println()
	fmt.Printf("  %s  v%s\n", primary.Render("gsh"), s.version)
	fmt.Printf("  %s\n", muted.Render("theme: "+s.cfg.Theme+"  |  type 'help' for built-ins"))
	fmt.Printf("  %s\n\n", secondary.Render("All commands run for real. Pipes, redirects, env vars all work."))
}

// parentDir returns the parent directory of a file path.
func parentDir(file string) string {
	if file == "" {
		return ""
	}
	idx := strings.LastIndex(file, "/")
	if idx <= 0 {
		return ""
	}
	return file[:idx]
}
