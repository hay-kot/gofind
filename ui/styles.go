package ui

import "github.com/charmbracelet/lipgloss"

const (
	ColorBlue   = lipgloss.Color("#255F85")
	ColorRed    = lipgloss.Color("#DA4167")
	ColorSubtle = lipgloss.Color("#848484")
	ColorWhite  = lipgloss.Color("#FFFFFF")
)

var (
	Bold         = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Render
	Subtle       = lipgloss.NewStyle().Foreground(ColorSubtle).Render
	AccentRed    = lipgloss.NewStyle().Foreground(ColorRed).Render
	AccentBlue   = lipgloss.NewStyle().Foreground(ColorBlue).Render
	HighlightRow = lipgloss.NewStyle().Background(lipgloss.Color("#2D2F27")).Render
)
