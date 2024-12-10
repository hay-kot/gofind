package tui

import (
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
