package app

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
	SuccessStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#16a34a", Dark: "#16a34a"}).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
	ErrorStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#b91c1c", Dark: "#b91c1c"}).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
	WarningStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.AdaptiveColor{Light: "#f97316", Dark: "#f97316"}).
			Margin(1, 1, 0, 0).
			Padding(0, 2).Render
)

var (
	audioQualityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	audioFormatStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 1).
				Render

	audioFileSizeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f97316")).
				Padding(0, 1).
				Render
)

var (
	videoQualityStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	videoFormatStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 1).
				Render

	videoResolutionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f97316")).
				Padding(0, 1).
				Render

	videoFileSizeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2563eb")).
				Padding(0, 1).
				Render
)

var (
	subtitleLangStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	subtitleSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#16a34a")).
				Padding(0, 1).
				Render
)
