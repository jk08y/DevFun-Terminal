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
	aliases       map[string]string
	dirStack      []string
}

// New creates a Shell from the provided configuration.
func New(cfg *config.Config, version string) *Shell {
	return &Shell{
		cfg:           cfg,
		version:       version,
		theme:         theme.Get(cfg.Theme),
		promptBuilder: prompt.New(cfg),
		aliases:       make(map[string]string),
		dirStack:      []string{},
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
	if dir := parentDir(s.cfg.History.File); dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

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

	runner, err := interp.New(
		interp.StdIO(os.Stdin, os.Stdout, os.Stderr),
		interp.Env(nil),
	)
	if err != nil {
		return fmt.Errorf("initialising shell runner: %w", err)
	}
	s.runner = runner

	// Load aliases and execute RC file before the first prompt
	s.loadAliases()
	if s.cfg.RCFile != "" {
		s.sourceFile(s.cfg.RCFile, true) // silent if missing
	}

	s.printWelcome()

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

		s.history = append(s.history, line)

		// Expand aliases then split with quote awareness
		expanded := s.expandAlias(line)
		args := splitFields(expanded)

		if result, ok := s.handleBuiltin(args); ok {
			s.exitCode = result.exitCode
			if result.doExit {
				break
			}
			continue
		}

		s.exitCode = s.execute(expanded)
	}

	return nil
}

// execute runs a shell command line and returns its exit code.
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

// sourceFile executes the commands in a file in the current shell context.
// If silent is true, a missing file is not treated as an error.
func (s *Shell) sourceFile(path string, silent bool) int {
	data, err := os.ReadFile(path)
	if err != nil {
		if silent {
			return 0
		}
		fmt.Fprintf(os.Stderr, "source: %v\n", err)
		return 1
	}

	f, err := syntax.NewParser().Parse(strings.NewReader(string(data)), path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "source: parse error in %s: %v\n", path, err)
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
