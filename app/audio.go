package app

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type AudioFormat struct {
	ID       string
	Format   string
	Quality  string
	Filesize string
}

type AudioFormatSelection struct {
	URL         string
	Formats     []AudioFormat
	Choice      int
	Selected    bool
	Downloading bool
	Done        bool
	Error       bool
	ErrMsg      string
}

func (m AppModel) fetchAudioFormats(url string) tea.Cmd {
	return func() tea.Msg {
		stdout, _, err := runYtdlp([]string{"-F", url})
		if err != nil {
			return AudioFormatMsg{Error: fmt.Sprintf("Error fetching formats: %v. Check output.log for details.", err)}
		}

		formats := ParseAudioFormats(stdout)

		return AudioFormatMsg{URL: url, Formats: formats}
	}
}

func ParseAudioFormats(output string) []AudioFormat {
	lines := strings.Split(output, "\n")
	var formats []AudioFormat

	for _, line := range lines {
		if strings.Contains(line, "Available formats") ||
			strings.Contains(line, "ID  EXT") ||
			strings.Contains(line, "----") ||
			strings.TrimSpace(line) == "" ||
			strings.Contains(line, "[youtube]") {
			continue
		}

		if strings.Contains(line, "audio only") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				id := fields[0]

				if strings.Contains(id, "-drc") {
					continue
				}

				quality := "Audio"
				filesize := "Unknown size"

				for _, field := range fields {
					if strings.HasSuffix(field, "k") {
						bitrateStr := strings.TrimSuffix(field, "k")
						if _, err := fmt.Sscanf(bitrateStr, "%f", new(float64)); err == nil {
							quality = bitrateStr + " kbps"
						}
					}

					if strings.Contains(field, "MiB") || strings.Contains(field, "KiB") {
						filesize = field
					}
				}

				if quality == "Audio" {
					if strings.Contains(line, "Default, high") {
						quality = "High quality"
					} else if strings.Contains(line, "Default, low") {
						quality = "Low quality"
					} else if strings.Contains(line, "[en]") {
						quality = "English audio"
					}
				}

				formatType := "audio"
				if len(fields) > 1 {
					ext := fields[1]
					if ext == "m4a" {
						formatType = "M4A (AAC)"
					} else if ext == "webm" {
						formatType = "WebM (Opus)"
					}
				}

				format := AudioFormat{
					ID:       id,
					Format:   formatType,
					Quality:  quality,
					Filesize: filesize,
				}

				exists := false
				for _, f := range formats {
					if f.ID == format.ID {
						exists = true
						break
					}
				}

				if !exists {
					formats = append(formats, format)
				}
			}
		}
	}

	sort.Slice(formats, func(i, j int) bool {

		iBitrate := extractBitrate(formats[i].Quality)
		jBitrate := extractBitrate(formats[j].Quality)

		if iBitrate == 0 && jBitrate == 0 {
			return formats[i].Quality > formats[j].Quality
		}

		return iBitrate > jBitrate
	})

	if len(formats) == 0 {
		formats = append(formats, AudioFormat{
			ID:       "bestaudio",
			Format:   "audio",
			Quality:  "Best quality",
			Filesize: "Unknown size",
		})
		formats = append(formats, AudioFormat{
			ID:       "worstaudio",
			Format:   "audio",
			Quality:  "Low quality",
			Filesize: "Unknown size",
		})
	}

	return formats
}

func extractBitrate(quality string) int {
	re := regexp.MustCompile(`(\d+)\s*kbps`)
	matches := re.FindStringSubmatch(quality)
	if len(matches) > 1 {
		if bitrate, err := strconv.Atoi(matches[1]); err == nil {
			return bitrate
		}
	}
	return 0
}

func (m AppModel) downloadAudio(url string, formatID string) tea.Cmd {
	return func() tea.Msg {
		outputTemplate := fmt.Sprintf("assets/%d_%%(title).120B [%%(id)s].%%(ext)s", time.Now().Unix())

		args := []string{"-f", formatID, "-x", "--audio-quality", "0"}
		args = append(args, throttleArgs...)
		args = append(args, "-o", outputTemplate, url)

		_, stderr, err := runYtdlp(args)
		if err != nil {

			if strings.Contains(stderr, "403") || strings.Contains(stderr, "Forbidden") {
				retry := []string{"-f", "bestaudio", "-x", "--audio-quality", "0"}
				retry = append(retry, throttleArgs...)
				retry = append(retry, "-o", outputTemplate, url)

				if _, _, err = runYtdlp(retry); err != nil {
					return AudioDownloadMsg{Error: fmt.Sprintf("Error downloading audio: %v. Check output.log for details.", err)}
				}
			} else {
				return AudioDownloadMsg{Error: fmt.Sprintf("Error downloading audio: %v. Check output.log for details.", err)}
			}
		}

		return AudioDownloadMsg{Done: true}
	}
}

type AudioFormatMsg struct {
	URL     string
	Formats []AudioFormat
	Error   string
}

type AudioDownloadMsg struct {
	Done  bool
	Error string
}
