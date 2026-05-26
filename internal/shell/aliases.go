package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// loadAliases reads the aliases file and populates s.aliases.
func (s *Shell) loadAliases() {
	if s.cfg.AliasFile == "" {
		return
	}
	f, err := os.Open(s.cfg.AliasFile)
	if err != nil {
		return // file doesn't exist yet; that's fine
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		s.aliases[name] = value
	}
}

// saveAliases writes the current aliases map to the aliases file.
func (s *Shell) saveAliases() {
	if s.cfg.AliasFile == "" {
		return
	}
	f, err := os.Create(s.cfg.AliasFile)
	if err != nil {
		return
	}
	defer f.Close()

	for name, value := range s.aliases {
		fmt.Fprintf(f, "%s=%s\n", name, value)
	}
}

// expandAlias replaces the first word of a command line with its alias expansion.
// It does not recurse to prevent infinite loops.
func (s *Shell) expandAlias(line string) string {
	if len(s.aliases) == 0 {
		return line
	}
	trimmed := strings.TrimSpace(line)
	// Find first word boundary
	idx := strings.IndexByte(trimmed, ' ')
	var cmd, rest string
	if idx < 0 {
		cmd = trimmed
		rest = ""
	} else {
		cmd = trimmed[:idx]
		rest = trimmed[idx:] // includes leading space
	}
	if expansion, ok := s.aliases[cmd]; ok {
		return expansion + rest
	}
	return line
}

// splitFields splits a string on whitespace while respecting single and double quotes.
// Quotes are stripped from the result tokens.
func splitFields(s string) []string {
	var fields []string
	var cur strings.Builder
	inSingle := false
	inDouble := false

	for _, r := range s {
		switch {
		case r == '\'' && !inDouble:
			inSingle = !inSingle
		case r == '"' && !inSingle:
			inDouble = !inDouble
		case (r == ' ' || r == '\t') && !inSingle && !inDouble:
			if cur.Len() > 0 {
				fields = append(fields, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		fields = append(fields, cur.String())
	}
	return fields
}
