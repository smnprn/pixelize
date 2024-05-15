package ui

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	BorderColor lipgloss.AdaptiveColor
	Framed      lipgloss.Style
	Success     lipgloss.Style
	Failure     lipgloss.Style
	Info        lipgloss.Style
}

func DefaultStyle() *Styles {
	s := new(Styles)

	s.BorderColor = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}

	s.Framed = lipgloss.
		NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(55)

	s.Success = lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}).
		Bold(true)

	s.Failure = lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}).
		Bold(true)

	s.Info = lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("#3C3C3C")).Align(lipgloss.Center)

	return s
}
