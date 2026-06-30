package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type VideoFormat struct {
	ID         string
	Format     string
	Quality    string
	Filesize   string
	Resolution string
}

type VideoFormatSelection struct {
	URL         string
	Formats     []VideoFormat
	Choice      int
	Selected    bool
	Downloading bool
	Done        bool
	Error       bool
	ErrMsg      string
}

func (m AppModel) fetchVideoFormats(url string) tea.Cmd {
	return func() tea.Msg {
		stdout, _, err := runYtdlp([]string{"-F", url})
		if err != nil {
			return VideoFormatMsg{Error: fmt.Sprintf("Error fetching formats: %v. Check output.log for details.", err)}
		}

		formats := ParseVideoFormats(stdout)

		return VideoFormatMsg{URL: url, Formats: formats}
	}
}

func ParseVideoFormats(output string) []VideoFormat {
	lines := strings.Split(output, "\n")
	var formats []VideoFormat

	for _, line := range lines {

		if strings.Contains(line, "Available formats") ||
			strings.Contains(line, "ID  EXT") ||
			strings.Contains(line, "----") ||
			strings.TrimSpace(line) == "" ||
			strings.Contains(line, "[youtube]") {
			continue
		}

		if strings.Contains(line, "audio only") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 3 {

			id := fields[0]

			format := "video"
			resolution := "Unknown resolution"
			filesize := "Unknown size"

			for _, field := range fields {
				if strings.Contains(field, "x") && strings.ContainsAny(field, "0123456789") {
					resolution = field
				}

				if strings.Contains(field, "MiB") || strings.Contains(field, "KiB") {
					filesize = field
				}
			}

			quality := resolution
			if strings.Contains(line, "360p") {
				quality = "360p"
			} else if strings.Contains(line, "480p") {
				quality = "480p"
			} else if strings.Contains(line, "720p") {
				quality = "720p HD"
			} else if strings.Contains(line, "1080p") {
				quality = "1080p Full HD"
			} else if strings.Contains(line, "1440p") {
				quality = "1440p Quad HD"
			} else if strings.Contains(line, "2160p") {
				quality = "2160p 4K"
			}

			videoFormat := VideoFormat{
				ID:         id,
				Format:     format,
				Quality:    quality,
				Filesize:   filesize,
				Resolution: resolution,
			}

			exists := false
			for _, f := range formats {
				if f.ID == videoFormat.ID {
					exists = true
					break
				}
			}

			if !exists {
				formats = append(formats, videoFormat)
			}
		}
	}

	if len(formats) == 0 {
		formats = append(formats, VideoFormat{
			ID:         "best",
			Format:     "video",
			Quality:    "Best quality",
			Filesize:   "Unknown size",
			Resolution: "Highest available",
		})
		formats = append(formats, VideoFormat{
			ID:         "worst",
			Format:     "video",
			Quality:    "Low quality",
			Filesize:   "Unknown size",
			Resolution: "Lowest available",
		})
	}

	return formats
}

func (m AppModel) downloadVideo(url string, formatID string) tea.Cmd {
	return func() tea.Msg {
		outputTemplate := fmt.Sprintf("assets/%d_%%(title).120B [%%(id)s].%%(ext)s", time.Now().Unix())

		args := []string{"-f", formatID}
		args = append(args, throttleArgs...)
		args = append(args, "-o", outputTemplate, url)

		_, stderr, err := runYtdlp(args)
		if err != nil {

			if strings.Contains(stderr, "403") || strings.Contains(stderr, "Forbidden") {
				retry := []string{"-f", "best"}
				retry = append(retry, throttleArgs...)
				retry = append(retry, "-o", outputTemplate, url)

				if _, _, err = runYtdlp(retry); err != nil {
					return VideoDownloadMsg{Error: fmt.Sprintf("Error downloading video: %v. Check output.log for details.", err)}
				}
			} else {
				return VideoDownloadMsg{Error: fmt.Sprintf("Error downloading video: %v. Check output.log for details.", err)}
			}
		}

		return VideoDownloadMsg{Done: true}
	}
}

type VideoFormatMsg struct {
	URL     string
	Formats []VideoFormat
	Error   string
}

type VideoDownloadMsg struct {
	Done  bool
	Error string
}
