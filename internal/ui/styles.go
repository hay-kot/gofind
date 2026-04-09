package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/gofind/internal/config"
)

var (
	ColorPrompt              lipgloss.Color
	ColorSelectedIndicator   lipgloss.Color
	ColorSecondaryText       lipgloss.Color
	ColorPrimaryText         lipgloss.Color
	ColorSelectionBackground lipgloss.Color
)

var (
	Bold         func(string) string
	Subtle       func(string) string
	AccentRed    func(string) string
	AccentBlue   func(string) string
	HighlightRow func(string) string
)

func Init(theme config.Theme) {
	ColorPrompt = lipgloss.Color(theme.Prompt)
	ColorSelectedIndicator = lipgloss.Color(theme.SelectedIndicator)
	ColorSecondaryText = lipgloss.Color(theme.SecondaryText)
	ColorPrimaryText = lipgloss.Color(theme.PrimaryText)
	ColorSelectionBackground = lipgloss.Color(theme.SelectionBackground)

	// Use stderr-based renderer so color detection works even when stdout is
	// piped (e.g. when invoked via shell substitution like $(repos)).
	// The TUI also outputs to stderr, so this matches its environment.
	r := lipgloss.NewRenderer(os.Stderr)

	boldStyle := r.NewStyle().Bold(true).Foreground(ColorPrimaryText)
	subtleStyle := r.NewStyle().Foreground(ColorSecondaryText)
	accentRedStyle := r.NewStyle().Foreground(ColorSelectedIndicator)
	accentBlueStyle := r.NewStyle().Foreground(ColorPrompt)
	highlightRowStyle := r.NewStyle().Background(ColorSelectionBackground)

	Bold = func(s string) string { return boldStyle.Render(s) }
	Subtle = func(s string) string { return subtleStyle.Render(s) }
	AccentRed = func(s string) string { return accentRedStyle.Render(s) }
	AccentBlue = func(s string) string { return accentBlueStyle.Render(s) }
	HighlightRow = func(s string) string { return highlightRowStyle.Render(s) }
}
