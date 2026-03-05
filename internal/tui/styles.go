package tui

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6600")).
		MarginBottom(1)

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CC00"))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAFF"))
)
