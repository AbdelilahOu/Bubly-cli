package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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

		os.MkdirAll("assets", 0755)

		var path, ffmpegPath string
		if isWindows() {
			path = "bin/yt-dlp.exe"
			ffmpegPath = "bin/ffmpeg.exe"
		} else {
			path = "bin/yt-dlp"
			ffmpegPath = "bin/ffmpeg"
		}

		logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return SubtitleLangMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		var outBuf, errBuf strings.Builder

		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil

		var args []string
		args = append(args, "--list-subs", url)

		if useFfmpeg {
			args = append(args, "--ffmpeg-location", ffmpegPath)
		} else {
			// Add a warning to the log file if ffmpeg is not found
			fmt.Fprintf(logFile, "Warning: ffmpeg not found. Some features may not work correctly.\n")
		}

		cmd := exec.Command(path, args...)
		cmd.Stdout = io.MultiWriter(&outBuf, logFile)
		cmd.Stderr = io.MultiWriter(&errBuf, logFile)

		err = cmd.Run()

		if err != nil {
			return SubtitleLangMsg{Error: fmt.Sprintf("Error fetching subtitle languages: %v. Check output.log for details.", err)}
		}

		languages := ParseSubtitleLanguages(outBuf.String())
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

		os.MkdirAll("assets", 0755)

		var path, ffmpegPath string
		if isWindows() {
			path = "bin/yt-dlp.exe"
			ffmpegPath = "bin/ffmpeg.exe"
		} else {
			path = "bin/yt-dlp"
			ffmpegPath = "bin/ffmpeg"
		}

		logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return SubtitleDownloadMsg{Error: fmt.Sprintf("Error creating log file: %v", err)}
		}
		defer logFile.Close()

		var outBuf, errBuf strings.Builder

		_, err = os.Stat(ffmpegPath)
		useFfmpeg := err == nil

		var args []string
		args = append(args, "--write-sub", "--write-auto-sub", "--sub-lang", langCode, "--skip-download")

		if useFfmpeg {
			args = append(args, "--ffmpeg-location", ffmpegPath)
		} else {
			// Add a warning to the log file if ffmpeg is not found
			fmt.Fprintf(logFile, "Warning: ffmpeg not found. Some features may not work correctly.\n")
		}

		args = append(args, "--sleep-requests", "1", "--sleep-interval", "5", "--max-sleep-interval", "10")
		outputTemplate := fmt.Sprintf("assets/%d_%%(title).120B [%%(id)s].%%(ext)s", time.Now().Unix())
		args = append(args, "-o", outputTemplate, url)

		cmd := exec.Command(path, args...)
		cmd.Stdout = io.MultiWriter(&outBuf, logFile)
		cmd.Stderr = io.MultiWriter(&errBuf, logFile)
		err = cmd.Run()

		if err != nil {
			errorOutput := errBuf.String()

			if strings.Contains(errorOutput, "429") || strings.Contains(errorOutput, "Too Many Requests") {
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
