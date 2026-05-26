package completer

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Completer implements readline.AutoCompleter for gsh.
// It completes executables (first token) and file paths (subsequent tokens).
type Completer struct{}

// New returns a ready-to-use Completer.
func New() *Completer {
	return &Completer{}
}

// Do is called by readline on every Tab press.
// line is the full line buffer; pos is the cursor position.
// Returns candidate suffixes and the length of the prefix already typed.
func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])
	tokens := strings.Fields(lineStr)

	// Determine what the user is currently typing
	var prefix string
	trailingSpace := len(lineStr) > 0 && lineStr[len(lineStr)-1] == ' '

	if len(tokens) == 0 || trailingSpace {
		prefix = ""
	} else {
		prefix = tokens[len(tokens)-1]
	}

	isCommand := len(tokens) == 0 || (len(tokens) == 1 && !trailingSpace)

	var candidates []string
	if isCommand {
		candidates = append(candidates, completeCommands(prefix)...)
		candidates = append(candidates, completeFiles(prefix)...)
	} else {
		candidates = completeFiles(prefix)
	}

	if len(candidates) == 0 {
		return nil, 0
	}

	result := make([][]rune, 0, len(candidates))
	for _, cand := range candidates {
		suffix := strings.TrimPrefix(cand, prefix)
		result = append(result, []rune(suffix))
	}

	return result, len([]rune(prefix))
}

// completeFiles returns file/dir names under the relevant directory that match prefix.
func completeFiles(prefix string) []string {
	dir := "."
	base := prefix

	if strings.Contains(prefix, string(filepath.Separator)) {
		if prefix[len(prefix)-1] == filepath.Separator {
			dir = prefix
			base = ""
		} else {
			dir = filepath.Dir(prefix)
			base = filepath.Base(prefix)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var out []string
	for _, e := range entries {
		name := e.Name()
		// Skip hidden unless user typed a dot
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(base, ".") {
			continue
		}
		if !strings.HasPrefix(name, base) {
			continue
		}
		full := filepath.Join(dir, name)
		if dir == "." {
			full = name
		}
		if e.IsDir() {
			full += string(filepath.Separator)
		}
		out = append(out, full)
	}
	return out
}

// completeCommands returns executable names on PATH that match prefix.
func completeCommands(prefix string) []string {
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, string(os.PathListSeparator))

	seen := make(map[string]bool)
	var out []string

	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return filepath.SkipDir
			}
			if d.IsDir() && path != dir {
				return filepath.SkipDir // don't recurse into subdirs of PATH entries
			}
			if d.IsDir() {
				return nil
			}
			name := filepath.Base(path)
			if !strings.HasPrefix(name, prefix) || seen[name] {
				return nil
			}
			info, err := d.Info()
			if err == nil && info.Mode()&0o111 != 0 {
				seen[name] = true
				out = append(out, name)
			}
			return nil
		})
	}
	return out
}
