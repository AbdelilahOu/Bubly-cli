package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SubtitleLanguage struct {
	Code string
	Name string
}

type SubtitleSelection struct {
	URL         string
	Languages   []SubtitleLanguage
	Choice      int
	Selected    bool
	Downloading bool
	Done        bool
	Error       bool
	ErrMsg      string
}

func (m AppModel) fetchSubtitleLanguages(url string) tea.Cmd {
	return func() tea.Msg {
		stdout, _, err := runYtdlp([]string{"--list-subs", url})
		if err != nil {
			return SubtitleLangMsg{Error: fmt.Sprintf("Error fetching subtitle languages: %v. Check output.log for details.", err)}
		}

		languages := ParseSubtitleLanguages(stdout)
		if len(languages) == 0 {
			return SubtitleLangMsg{Error: "No subtitles were found for this video."}
		}

		return SubtitleLangMsg{URL: url, Languages: languages}
	}
}

func ParseSubtitleLanguages(output string) []SubtitleLanguage {
	lines := strings.Split(output, "\n")
	var languages []SubtitleLanguage

	parsingSubtitles := false

	for _, line := range lines {

		if strings.Contains(line, "Available subtitles for") ||
			strings.Contains(line, "Available automatic captions for") {
			parsingSubtitles = true
			continue
		}

		if !parsingSubtitles ||
			strings.Contains(line, "Language Name") ||
			strings.Contains(line, "----") ||
			strings.TrimSpace(line) == "" ||
			strings.Contains(line, "[youtube]") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 1 {

			code := fields[0]

			if code == "en-orig" {
				continue
			}

			name := code
			switch code {
			case "en":
				name = "English"
			case "es":
				name = "Spanish"
			case "fr":
				name = "French"
			case "de":
				name = "German"
			case "it":
				name = "Italian"
			case "pt":
				name = "Portuguese"
			case "ru":
				name = "Russian"
			case "ja":
				name = "Japanese"
			case "ko":
				name = "Korean"
			case "zh-Hans":
				name = "Chinese (Simplified)"
			case "zh-Hant":
				name = "Chinese (Traditional)"
			case "ar":
				name = "Arabic"
			case "hi":
				name = "Hindi"
			case "tr":
				name = "Turkish"
			default:

				if len(name) > 0 {
					name = strings.ToUpper(name[:1]) + name[1:]
				}
			}

			lang := SubtitleLanguage{
				Code: code,
				Name: name,
			}

			exists := false
			for _, l := range languages {
				if l.Code == lang.Code {
					exists = true
					break
				}
			}

			if !exists {
				languages = append(languages, lang)
			}
		}
	}

	return languages
}

func (m AppModel) downloadSubtitles(url string, langCode string) tea.Cmd {
	return func() tea.Msg {
		outputTemplate := fmt.Sprintf("assets/%d_%%(title).120B [%%(id)s].%%(ext)s", time.Now().Unix())

		args := []string{"--write-sub", "--write-auto-sub", "--sub-lang", langCode, "--skip-download"}
		args = append(args, throttleArgs...)
		args = append(args, "-o", outputTemplate, url)

		_, stderr, err := runYtdlp(args)
		if err != nil {
			if strings.Contains(stderr, "429") || strings.Contains(stderr, "Too Many Requests") {
				return SubtitleDownloadMsg{Error: "Rate limited by YouTube. Please try again later."}
			}
			return SubtitleDownloadMsg{Error: fmt.Sprintf("Error downloading subtitles: %v. Check output.log for details.", err)}
		}

		return SubtitleDownloadMsg{Done: true}
	}
}

type SubtitleLangMsg struct {
	URL       string
	Languages []SubtitleLanguage
	Error     string
}

type SubtitleDownloadMsg struct {
	Done  bool
	Error string
}
