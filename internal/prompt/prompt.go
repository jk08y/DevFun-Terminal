package prompt

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jk08y/nexterm/internal/config"
	"github.com/jk08y/nexterm/internal/theme"
)

// Builder constructs the shell prompt string.
type Builder struct {
	cfg   *config.Config
	theme theme.Theme
}

// New creates a Builder using the configured theme.
func New(cfg *config.Config) *Builder {
	return &Builder{
		cfg:   cfg,
		theme: theme.Get(cfg.Theme),
	}
}

// SetTheme swaps the active theme at runtime.
func (b *Builder) SetTheme(t theme.Theme) {
	b.theme = t
}

// Build returns the rendered prompt string.
// exitCode is the exit status of the previous command (0 = success).
func (b *Builder) Build(exitCode int) string {
	primary := lipgloss.NewStyle().Foreground(b.theme.Primary).Bold(true)
	secondary := lipgloss.NewStyle().Foreground(b.theme.Secondary).Bold(true)
	muted := lipgloss.NewStyle().Foreground(b.theme.Muted)
	warnStyle := lipgloss.NewStyle().Foreground(b.theme.Warning)

	var sb strings.Builder

	// ── user[@host] ──────────────────────────────────────────────────────────
	if b.cfg.ShowUser {
		user := os.Getenv("USER")
		if user == "" {
			user = "user"
		}
		sb.WriteString(primary.Render(user))

		if b.cfg.ShowHost {
			host, _ := os.Hostname()
			sb.WriteString(muted.Render("@"))
			sb.WriteString(primary.Render(host))
		}
		sb.WriteString(" ")
	}

	// ── current directory ────────────────────────────────────────────────────
	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(cwd, home) {
		cwd = "~" + cwd[len(home):]
	}
	// Truncate deep paths: show only last 2 segments
	parts := strings.Split(filepath.ToSlash(cwd), "/")
	if len(parts) > 3 {
		cwd = "…/" + strings.Join(parts[len(parts)-2:], "/")
	}
	sb.WriteString(secondary.Render(cwd))

	// ── git branch ───────────────────────────────────────────────────────────
	if b.cfg.ShowGit {
		if branch := gitBranch(); branch != "" {
			sb.WriteString(" ")
			sb.WriteString(warnStyle.Render(" " + branch))
		}
	}

	// ── arrow (colour indicates last exit status) ────────────────────────────
	sb.WriteString("\n")
	arrowColor := b.theme.Success
	if exitCode != 0 {
		arrowColor = b.theme.Error
	}
	arrow := lipgloss.NewStyle().Foreground(arrowColor).Bold(true).Render("❯")
	sb.WriteString(arrow)
	sb.WriteString(" ")

	return sb.String()
}

// gitBranch returns the current git branch name, or "" if not in a repo.
func gitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	branch := strings.TrimSpace(string(out))
	if branch == "HEAD" {
		// Detached HEAD — show short hash instead
		cmd2 := exec.Command("git", "rev-parse", "--short", "HEAD")
		cmd2.Stderr = nil
		out2, err2 := cmd2.Output()
		if err2 == nil {
			return "(" + strings.TrimSpace(string(out2)) + ")"
		}
		return "(detached)"
	}
	return branch
}
