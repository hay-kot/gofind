package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jwalton/go-supportscolor"
	"github.com/muesli/termenv"
)

func setup() {
	term := supportscolor.Stderr()
	if term.Has16m {
		lipgloss.SetColorProfile(termenv.TrueColor)
	} else if term.Has256 {
		lipgloss.SetColorProfile(termenv.ANSI256)
	} else {
		lipgloss.SetColorProfile(termenv.ANSI)
	}

}

// Copy-pasted from github.com/muesli/termenv@v0.9.0/termenv_unix.go.
// TODO: Refactor after, [feature](https://ï.at/stderr) implemented.
func colorProfile() termenv.Profile {
	term := os.Getenv("TERM")
	colorTerm := os.Getenv("COLORTERM")

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if term == "screen" || !strings.HasPrefix(term, "screen") {
			// enable TrueColor in tmux, but not for old-school screen
			return termenv.TrueColor
		}
	case "yes":
		fallthrough
	case "true":
		return termenv.ANSI256
	}

	if strings.Contains(term, "256color") {
		return termenv.ANSI256
	}
	if strings.Contains(term, "color") {
		return termenv.ANSI
	}

	return termenv.Ascii
}
