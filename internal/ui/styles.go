package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/gofind/internal/core/config"
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

	boldStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorPrimaryText)
	subtleStyle := lipgloss.NewStyle().Foreground(ColorSecondaryText)
	accentRedStyle := lipgloss.NewStyle().Foreground(ColorSelectedIndicator)
	accentBlueStyle := lipgloss.NewStyle().Foreground(ColorPrompt)
	highlightRowStyle := lipgloss.NewStyle().Background(ColorSelectionBackground)

	Bold = func(s string) string { return boldStyle.Render(s) }
	Subtle = func(s string) string { return subtleStyle.Render(s) }
	AccentRed = func(s string) string { return accentRedStyle.Render(s) }
	AccentBlue = func(s string) string { return accentBlueStyle.Render(s) }
	HighlightRow = func(s string) string { return highlightRowStyle.Render(s) }
}
