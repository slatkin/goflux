package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("62") // Purpleish
	ColorSecondary = lipgloss.Color("230")
	ColorText      = lipgloss.Color("252")
	ColorDim       = lipgloss.Color("240")
	ColorError     = lipgloss.Color("196")
	ColorSuccess   = lipgloss.Color("46")

	// Styles
	StyleBase = lipgloss.NewStyle().Foreground(ColorText)

	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Background(ColorPrimary).
			Bold(true)

	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	StyleStatusUnread = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")). // Bright white
				Bold(true)

	StyleStatusRead = lipgloss.NewStyle().
			Foreground(ColorDim)

	StyleErrorMessage = lipgloss.NewStyle().
				Foreground(ColorError).
				Bold(true)
)
