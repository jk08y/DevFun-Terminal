package theme

import "github.com/charmbracelet/lipgloss"

// Theme holds a named colour palette used across the prompt and output.
type Theme struct {
	Name      string
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Muted     lipgloss.Color
	Text      lipgloss.Color
}

var themes = map[string]Theme{
	"dracula": {
		Name:      "dracula",
		Primary:   "#bd93f9",
		Secondary: "#ff79c6",
		Success:   "#50fa7b",
		Warning:   "#f1fa8c",
		Error:     "#ff5555",
		Muted:     "#6272a4",
		Text:      "#f8f8f2",
	},
	"nord": {
		Name:      "nord",
		Primary:   "#88c0d0",
		Secondary: "#81a1c1",
		Success:   "#a3be8c",
		Warning:   "#ebcb8b",
		Error:     "#bf616a",
		Muted:     "#4c566a",
		Text:      "#eceff4",
	},
	"catppuccin": {
		Name:      "catppuccin",
		Primary:   "#cba6f7",
		Secondary: "#f38ba8",
		Success:   "#a6e3a1",
		Warning:   "#f9e2af",
		Error:     "#f38ba8",
		Muted:     "#585b70",
		Text:      "#cdd6f4",
	},
	"onedark": {
		Name:      "onedark",
		Primary:   "#61afef",
		Secondary: "#c678dd",
		Success:   "#98c379",
		Warning:   "#e5c07b",
		Error:     "#e06c75",
		Muted:     "#5c6370",
		Text:      "#abb2bf",
	},
	"tokyo-night": {
		Name:      "tokyo-night",
		Primary:   "#7aa2f7",
		Secondary: "#bb9af7",
		Success:   "#9ece6a",
		Warning:   "#e0af68",
		Error:     "#f7768e",
		Muted:     "#565f89",
		Text:      "#c0caf5",
	},
}

// Get returns the named theme, falling back to dracula if unknown.
func Get(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["dracula"]
}

// List returns all available theme names.
func List() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}
