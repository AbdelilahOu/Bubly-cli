package app

import (
	"github.com/AbdelilahOu/Bubly-cli-app/types"
	"github.com/AbdelilahOu/Bubly-cli-app/utils"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewsOptions struct {
	View        string
	ChoiceLabel string
}

type AppModel struct {
	Choice               int
	Quitting             bool
	History              []string
	Textarea             textarea.Model
	Text                 string
	IsTextAreaActive     bool
	IsUrlWritten         bool
	PrintingIsDone       bool
	PrintingError        bool
	Warning              string
	CheckingYtdlp        bool
	InstallingYtdlp      bool
	InstallationProgress int
	InstallationTotal    int
	InstallationMessage  string
	YtdlpInstalled       bool
	AudioFormatSel       *AudioFormatSelection
	VideoFormatSel       *VideoFormatSelection
	SubtitleSel          *SubtitleSelection
	Page                 int
	ItemsPerPage         int
	Width                int
	Height               int
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return types.CheckYtdlpMsg{Installed: utils.CheckYtdlp()}
		},
	)
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		textareaWidth := viewWidth(m) - 24
		if textareaWidth < 36 {
			textareaWidth = 36
		}
		m.Textarea.SetWidth(textareaWidth)
		return m, nil
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if m.IsTextAreaActive {
			if k == "esc" || k == "ctrl+c" {
				m.Quitting = true
				return m, tea.Quit
			}
		} else {
			if k == "q" || k == "esc" || k == "ctrl+c" {
				m.Quitting = true
				return m, tea.Quit
			}
		}
		if k == "backspace" && len(m.History) > 0 {
			if m.Textarea.Value() == "" {
				m.IsUrlWritten = false
				m.PrintingError = false
				m.PrintingIsDone = false
				m = removeFromHistory(m)
			} else {
				m.Text = m.Textarea.Value()[:len(m.Textarea.Value())-1]
			}
		}
	}

	if m.CheckingYtdlp || m.InstallingYtdlp {
		return UpdateYtdlp(msg, m)
	}

	return UpdateYoutube(msg, m)
}

func (m AppModel) View() string {
	if m.Quitting {
		return renderShell(m, SuccessStyle("See you later! 👋"))
	}

	if m.CheckingYtdlp || m.InstallingYtdlp {
		return renderShell(m, YtdlpView(m))
	}

	content := YoutubeView(m)
	if m.Warning != "" {
		content += "\n\n" + WarningStyle(m.Warning)
	}
	return renderShell(m, content)
}

func UpdateYtdlp(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.CheckYtdlpMsg:
		m.YtdlpInstalled = msg.Installed
		if !m.YtdlpInstalled {
			m.CheckingYtdlp = true
		} else {
			m.CheckingYtdlp = false
		}
		return m, nil
	case types.YtdlpInstalledMsg:
		if msg.Err != nil {
			m.Warning = "Error installing yt-dlp: " + msg.Err.Error()
		} else {
			m.Warning = "yt-dlp installed successfully"
			m.YtdlpInstalled = true
		}
		m.CheckingYtdlp = false
		m.InstallingYtdlp = false
		return m, nil
	case types.ProgressMsg:
		m.InstallationProgress = msg.Progress
		m.InstallationTotal = msg.Total
		m.InstallationMessage = msg.Message
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.Choice == 0 {
				if m.CheckingYtdlp {
					m.InstallingYtdlp = true
					return m, utils.InstallYtdlp()
				}
			} else {
				if m.CheckingYtdlp {
					m.CheckingYtdlp = false
					m.Warning = "yt-dlp is not installed. Some features may not work."
				}
				return m, nil
			}
		case "up", "k":
			if m.Choice > 0 {
				m.Choice--
			}
		case "down", "j":
			if m.Choice < 1 {
				m.Choice++
			}
		}
	}
	return m, nil
}

func YtdlpView(m AppModel) string {
	var s string

	if m.CheckingYtdlp {
		if m.InstallingYtdlp {
			s += TitleStyle("Environment Setup") + "\n\n"
			s += "Installing `yt-dlp` in local `bin/`...\n"
			s += "This only needs to happen once.\n"
		} else {
			s += TitleStyle("Environment Setup") + "\n\n"
			s += "yt-dlp was not found in local `bin/`.\n"
			s += "Install it now?\n\n"
		}
	} else {
		return ""
	}

	if !m.InstallingYtdlp {
		choices := []string{"Yes", "No"}
		for i, choice := range choices {
			s += checkbox(choice, m.Choice == i) + "\n"
		}
	}

	return s
}

func destructureOptions(options []ViewsOptions, c int) []any {
	var choices []any
	for i, option := range options {
		choices = append(choices, checkbox(option.ChoiceLabel, c == i))
	}
	return choices
}

func checkbox(label string, checked bool) string {
	if checked {
		return selectedOptionStyle.Render("▶ " + label)
	}
	return optionStyle.Render("  " + label)
}

func appendToHistory(m AppModel, s string) AppModel {
	m.History = append(m.History, s)
	return m
}
func removeFromHistory(m AppModel) AppModel {

	if len(m.History) > 0 {
		switch m.History[len(m.History)-1] {
		case "yt-download-audio":

			m.AudioFormatSel = nil
			m.IsUrlWritten = false
			m.Text = ""
			m.Textarea.Reset()
		case "yt-download-video":

			m.VideoFormatSel = nil
			m.IsUrlWritten = false
			m.Text = ""
			m.Textarea.Reset()
		case "yt-download-subtitles":

			m.SubtitleSel = nil
			m.IsUrlWritten = false
			m.Text = ""
			m.Textarea.Reset()
		}
	}

	m.History = m.History[:len(m.History)-1]
	m.Choice = 0

	m.PrintingError = false
	m.PrintingIsDone = false
	m.Warning = ""

	return m
}
