package main

import (
	"fmt"

	"github.com/AbdelilahOu/Bubly-cli-app/app"
	"github.com/AbdelilahOu/Bubly-cli-app/utils"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	utils.ClearTerminal()

	ta := textarea.New()
	ta.Placeholder = "Paste a YouTube URL and press Enter..."
	ta.Focus()

	ta.Prompt = "▶ "
	ta.CharLimit = 280

	ta.SetWidth(72)
	ta.SetHeight(2)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#111827"))
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ca3af"))
	ta.BlurredStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280"))
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ca3af"))

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	initialModel := app.AppModel{
		Choice:           0,
		Quitting:         false,
		History:          []string{},
		Textarea:         ta,
		Text:             "",
		IsTextAreaActive: false,
		IsUrlWritten:     false,
		PrintingIsDone:   false,
		PrintingError:    false,
		CheckingYtdlp:    true,
		Page:             0,
		ItemsPerPage:     5,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
