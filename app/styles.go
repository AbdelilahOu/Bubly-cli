package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	pageStyle = lipgloss.NewStyle().
			Padding(1, 2)

	shellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	headerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 2)

	headerSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6b7280")).
				PaddingLeft(1)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 2)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#d1d5db")).
			Padding(1, 2)

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#374151")).
			Padding(0, 1)

	selectedOptionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#16a34a")).
			Padding(0, 2).Render

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#b91c1c")).
			Padding(0, 2).Render

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#f97316")).
			Padding(0, 2).Render

	TitleStyle = sectionTitleStyle.Render

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

	subtitleLangStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Padding(0, 1).
				Render

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280"))
)

func viewWidth(m AppModel) int {
	if m.Width <= 0 {
		return 96
	}
	w := m.Width - 6
	if w < 72 {
		return 72
	}
	if w > 130 {
		return 130
	}
	return w
}

func renderShell(m AppModel, body string) string {
	width := viewWidth(m)
	title := headerTitleStyle.Render(" Bubly CLI ")
	header := lipgloss.JoinHorizontal(lipgloss.Left, title)

	status := "Ready"
	if m.CheckingYtdlp {
		status = "Setup"
	}
	if m.Warning != "" {
		status = "Warning"
	}
	info := helpStyle.Render(fmt.Sprintf("status: %s  •  page: %d", status, m.Page+1))
	helpLine := helpStyle.Render("j/k: move • enter: select • backspace: back • q: quit")

	content := cardStyle.Width(width - 6).Render(strings.TrimRight(body, "\n"))
	shell := shellStyle.Width(width).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			content,
			"",
			info,
			helpLine,
		),
	)
	return pageStyle.Render(shell)
}
