package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type spinnerTickMsg time.Time

func spinnerTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

var YoutubeOptions = []ViewsOptions{
	{
		View:        "yt-download-video",
		ChoiceLabel: "Download Youtube video 📥",
	},
	{
		View:        "yt-download-audio",
		ChoiceLabel: "Download Youtube audio 🎵",
	},
	{
		View:        "yt-download-subtitles",
		ChoiceLabel: "Download Youtube subtitles 📝",
	},
}

func UpdateYoutube(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {

	if len(m.History) > 0 && m.History[0] == "yt-download-audio" && m.IsUrlWritten && m.AudioFormatSel != nil {
		sel := m.AudioFormatSel
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if c, p, ok := moveSelection(msg.String(), sel.Choice, len(sel.Formats), m.Page, m.ItemsPerPage); ok {
				sel.Choice, m.Page = c, p
				return m, nil
			}
			if msg.String() == "enter" && !sel.Selected {
				sel.Selected = true
				sel.Downloading = true
				formatID := sel.Formats[sel.Choice].ID
				return m, tea.Batch(m.downloadAudio(sel.URL, formatID), spinnerTick())
			}
			return m, nil
		case AudioDownloadMsg:
			sel.Downloading = false
			if msg.Error != "" {
				sel.Error = true
				sel.ErrMsg = msg.Error
			} else {
				sel.Done = true
			}
			return m, nil
		case spinnerTickMsg:
			if sel.Downloading {
				return m, spinnerTick()
			}
			return m, nil
		}
		return m, nil
	}

	if len(m.History) > 0 && m.History[0] == "yt-download-video" && m.IsUrlWritten && m.VideoFormatSel != nil {
		sel := m.VideoFormatSel
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if c, p, ok := moveSelection(msg.String(), sel.Choice, len(sel.Formats), m.Page, m.ItemsPerPage); ok {
				sel.Choice, m.Page = c, p
				return m, nil
			}
			if msg.String() == "enter" && !sel.Selected {
				sel.Selected = true
				sel.Downloading = true
				formatID := sel.Formats[sel.Choice].ID
				return m, tea.Batch(m.downloadVideo(sel.URL, formatID), spinnerTick())
			}
			return m, nil
		case VideoDownloadMsg:
			sel.Downloading = false
			if msg.Error != "" {
				sel.Error = true
				sel.ErrMsg = msg.Error
			} else {
				sel.Done = true
			}
			return m, nil
		case spinnerTickMsg:
			if sel.Downloading {
				return m, spinnerTick()
			}
			return m, nil
		}
		return m, nil
	}

	if len(m.History) > 0 && m.History[0] == "yt-download-subtitles" && m.IsUrlWritten && m.SubtitleSel != nil {
		sel := m.SubtitleSel
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if c, p, ok := moveSelection(msg.String(), sel.Choice, len(sel.Languages), m.Page, m.ItemsPerPage); ok {
				sel.Choice, m.Page = c, p
				return m, nil
			}
			if msg.String() == "enter" && !sel.Selected {
				sel.Selected = true
				sel.Downloading = true
				langCode := sel.Languages[sel.Choice].Code
				return m, tea.Batch(m.downloadSubtitles(sel.URL, langCode), spinnerTick())
			}
			return m, nil
		case SubtitleDownloadMsg:
			sel.Downloading = false
			if msg.Error != "" {
				sel.Error = true
				sel.ErrMsg = msg.Error
			} else {
				sel.Done = true
			}
			return m, nil
		case spinnerTickMsg:
			if sel.Downloading {
				return m, spinnerTick()
			}
			return m, nil
		}
		return m, nil
	}

	if len(m.History) > 0 {
		switch m.History[0] {
		case "yt-download-video":
			return UpdateDownloadVideo(msg, m)
		case "yt-download-audio":
			return UpdateDownloadAudio(msg, m)
		case "yt-download-subtitles":
			return UpdateDownloadSubtitles(msg, m)
		}
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(YoutubeOptions) > m.Choice+1 {
				m.Choice++
			}
		case "k", "up":
			if m.Choice > 0 {
				m.Choice--
			}
		case "enter":
			m.IsTextAreaActive = true
			m = appendToHistory(m, YoutubeOptions[m.Choice].View)
			return m, nil
		}
	case AudioFormatMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
		} else {
			m.AudioFormatSel = &AudioFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
			m.Page = 0
		}
		return m, nil
	case VideoFormatMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
		} else {
			m.VideoFormatSel = &VideoFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
			m.Page = 0
		}
		return m, nil
	case SubtitleLangMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
		} else {
			m.SubtitleSel = &SubtitleSelection{
				URL:       msg.URL,
				Languages: msg.Languages,
				Choice:    0,
			}
			m.Page = 0
		}
		return m, nil
	case AudioDownloadMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		}
		return m, nil
	case VideoDownloadMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		}
		return m, nil
	case SubtitleDownloadMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		}
		return m, nil
	}
	return m, nil
}

func YoutubeView(m AppModel) string {
	c := m.Choice
	var s strings.Builder

	if len(m.History) > 0 {
		switch m.History[0] {
		case "yt-download-video":
			s.WriteString(DownloadVideoView(m))
		case "yt-download-audio":
			s.WriteString(DownloadAudioView(m))
		case "yt-download-subtitles":
			s.WriteString(DownloadSubtitlesView(m))
		}
	} else {
		s.WriteString(TitleStyle("Choose Your Workflow"))
		s.WriteString("\n\n")
		s.WriteString("Pick a mode to start downloading media.\n\n")

		choices := fmt.Sprintf(
			strings.Repeat("%s\n", len(YoutubeOptions)),
			destructureOptions(YoutubeOptions, c)...,
		)
		s.WriteString(strings.TrimRight(choices, "\n"))
	}

	return s.String()
}

func DownloadVideoView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube video 📥"))
	s.WriteString("\n\n")

	if !m.IsUrlWritten {
		s.WriteString(m.Textarea.View())
		return s.String()
	}

	sel := m.VideoFormatSel
	if sel == nil {
		s.WriteString("Fetching available video formats for: " + m.Text + "\n")
		return s.String()
	}

	switch {
	case sel.Error:
		s.WriteString(ErrorStyle("Error: " + sel.ErrMsg))
	case sel.Done:
		s.WriteString(SuccessStyle("Video downloaded successfully! Check assets folder"))
	case sel.Downloading:
		s.WriteString("📥 Downloading video " + downloadSpinner() + "\n")
		s.WriteString("This may take a few moments...")
	case len(sel.Formats) > 0:
		s.WriteString("Select video format:\n\n")
		lines := make([]string, len(sel.Formats))
		for i, format := range sel.Formats {
			lines[i] = fmt.Sprintf("%s %s %s %s",
				videoQualityStyle(fixedCol(format.Quality, 14)),
				videoFormatStyle(fixedCol(format.Format, 10)),
				videoResolutionStyle(fixedCol(format.Resolution, 16)),
				videoFileSizeStyle(fixedCol(format.Filesize, 10)))
		}
		renderSelectList(&s, lines, sel.Choice, m.Page, m.ItemsPerPage)
		s.WriteString("\n\n(Press ↑/↓ to select, Enter to download, h/l for pagination)")
	default:
		s.WriteString("Loading available formats...")
	}
	return s.String()
}

func UpdateDownloadVideo(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)
	m.Textarea, tiCmd = m.Textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.IsUrlWritten {
				url := strings.TrimSpace(m.Textarea.Value())
				if url == "" {
					return m, nil
				}
				m.Text = url
				m.Textarea.Reset()
				m.IsUrlWritten = true
				m.IsTextAreaActive = false

				return m, m.fetchVideoFormats(m.Text)
			}
			return m, nil
		}
	case VideoFormatMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		} else {
			m.VideoFormatSel = &VideoFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
		}
		return m, nil
	case VideoDownloadMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		}
		return m, nil
	}
	return m, tea.Batch(tiCmd)
}

func DownloadAudioView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube audio \U0001F3B5"))
	s.WriteString("\n\n")

	if !m.IsUrlWritten {
		s.WriteString(m.Textarea.View())
		return s.String()
	}

	sel := m.AudioFormatSel
	if sel == nil {
		s.WriteString("Fetching available audio formats for: " + m.Text + "\n")
		return s.String()
	}

	switch {
	case sel.Error:
		s.WriteString(ErrorStyle("Error: " + sel.ErrMsg))
	case sel.Done:
		s.WriteString(SuccessStyle("Audio downloaded successfully! Check assets folder"))
	case sel.Downloading:
		s.WriteString("🔊 Downloading audio " + downloadSpinner() + "\n")
		s.WriteString("This may take a few moments...")
	case len(sel.Formats) > 0:
		s.WriteString("Select audio format:\n\n")
		lines := make([]string, len(sel.Formats))
		for i, format := range sel.Formats {
			lines[i] = fmt.Sprintf("%s %s %s",
				audioQualityStyle(fixedCol(format.Quality, 14)),
				audioFormatStyle(fixedCol(format.Format, 12)),
				audioFileSizeStyle(fixedCol(format.Filesize, 10)))
		}
		renderSelectList(&s, lines, sel.Choice, m.Page, m.ItemsPerPage)
		s.WriteString("\n\n(Press ↑/↓ to select, Enter to download, h/l for pagination)")
	default:
		s.WriteString("Loading available formats...")
	}
	return s.String()
}

func UpdateDownloadAudio(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)
	m.Textarea, tiCmd = m.Textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.IsUrlWritten {
				url := strings.TrimSpace(m.Textarea.Value())
				if url == "" {
					return m, nil
				}
				m.Text = url
				m.Textarea.Reset()
				m.IsUrlWritten = true
				m.IsTextAreaActive = false
				return m, m.fetchAudioFormats(m.Text)
			}
			return m, nil
		}
	case AudioFormatMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
			m.IsUrlWritten = false
		} else {
			m.AudioFormatSel = &AudioFormatSelection{
				URL:     msg.URL,
				Formats: msg.Formats,
				Choice:  0,
			}
		}
		return m, nil
	}
	return m, tea.Batch(tiCmd)
}

func DownloadSubtitlesView(m AppModel) string {
	var s strings.Builder
	s.WriteString(TitleStyle("Download Youtube subtitles \U0001F4DD"))
	s.WriteString("\n\n")

	if !m.IsUrlWritten {
		s.WriteString(m.Textarea.View())
		return s.String()
	}

	sel := m.SubtitleSel
	if sel == nil {
		s.WriteString("Fetching available subtitle languages for: " + m.Text + "\n")
		return s.String()
	}

	switch {
	case sel.Error:
		s.WriteString(ErrorStyle("Error: " + sel.ErrMsg))
	case sel.Done:
		s.WriteString(SuccessStyle("Subtitles downloaded successfully! Check assets folder"))
	case sel.Downloading:
		selectedLang := sel.Languages[sel.Choice].Name
		s.WriteString("📝 Downloading " + selectedLang + " subtitles " + downloadSpinner() + "\n")
		s.WriteString("This may take a few moments...")
	case len(sel.Languages) > 0:
		s.WriteString("Select subtitle language:\n\n")
		lines := make([]string, len(sel.Languages))
		for i, lang := range sel.Languages {
			lines[i] = subtitleLangStyle(lang.Name)
		}
		renderSelectList(&s, lines, sel.Choice, m.Page, m.ItemsPerPage)
		s.WriteString("\n\n(Press ↑/↓ to select, Enter to download, h/l for pagination)")
	default:
		s.WriteString("Loading available languages...")
	}
	return s.String()
}

func UpdateDownloadSubtitles(msg tea.Msg, m AppModel) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)
	m.Textarea, tiCmd = m.Textarea.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.IsUrlWritten {
				url := strings.TrimSpace(m.Textarea.Value())
				if url == "" {
					return m, nil
				}
				m.Text = url
				m.Textarea.Reset()
				m.IsUrlWritten = true
				m.IsTextAreaActive = false

				return m, m.fetchSubtitleLanguages(m.Text)
			}
			return m, nil
		}
	case SubtitleLangMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		} else {
			m.SubtitleSel = &SubtitleSelection{
				URL:       msg.URL,
				Languages: msg.Languages,
				Choice:    0,
			}
		}
		return m, nil
	case SubtitleDownloadMsg:
		if msg.Error != "" {
			m.Warning = msg.Error
		}
		return m, nil
	}
	return m, tea.Batch(tiCmd)
}

func moveSelection(key string, choice, count, page, perPage int) (newChoice, newPage int, ok bool) {
	if perPage < 1 {
		perPage = 1
	}
	switch key {
	case "j", "down":
		if choice+1 < count {
			choice++
		}
	case "k", "up":
		if choice > 0 {
			choice--
		}
	case "h", "left":
		if page > 0 {
			page--
			choice = page * perPage
		}
		return choice, page, true
	case "l", "right":
		totalPages := (count + perPage - 1) / perPage
		if page < totalPages-1 {
			page++
			choice = page * perPage
		}
		return choice, page, true
	default:
		return choice, page, false
	}
	return choice, choice / perPage, true
}

func renderSelectList(s *strings.Builder, lines []string, choice, page, perPage int) {
	if perPage < 1 {
		perPage = 1
	}
	total := len(lines)
	totalPages := (total + perPage - 1) / perPage
	if page >= totalPages {
		page = totalPages - 1
	}
	if page < 0 {
		page = 0
	}

	start := page * perPage
	end := start + perPage
	if end > total {
		end = total
	}

	for i := start; i < end; i++ {
		cursor := "  "
		style := optionStyle
		if i == choice {
			cursor = "> "
			style = selectedOptionStyle
		}
		s.WriteString(style.Render(cursor+lines[i]) + "\n")
	}

	if totalPages > 1 {
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("Page %d of %d | ", page+1, totalPages))
		if page > 0 {
			s.WriteString("<-- Previous (h) ")
		}
		if page < totalPages-1 {
			s.WriteString("Next (l) -->")
		}
	}
}

func downloadSpinner() string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := time.Now().UnixNano() / 100000000 % int64(len(frames))
	return frames[frame]
}

func fixedCol(value string, width int) string {
	if width <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) > width {
		if width == 1 {
			return "…"
		}
		return string(runes[:width-1]) + "…"
	}
	if len(runes) < width {
		return value + strings.Repeat(" ", width-len(runes))
	}
	return value
}
