package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/AbdelilahOu/Bubly-cli-app/types"
	tea "github.com/charmbracelet/bubbletea"
)

func CheckYtdlp() bool {
	var path string
	if runtime.GOOS == "windows" {
		path = "bin/yt-dlp.exe"
	} else {
		path = "bin/yt-dlp"
	}
	_, err := os.Stat(path)
	return err == nil
}

func InstallYtdlp() tea.Cmd {
	return func() tea.Msg {
		err := doInstallYtdlp()
		return types.YtdlpInstalledMsg{Err: err}
	}
}

func doInstallYtdlp() error {
	err := os.MkdirAll("bin", 0755)
	if err != nil {
		return err
	}

	var url string
	switch runtime.GOOS {
	case "windows":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	case "linux":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
	case "darwin":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
	default:
		return fmt.Errorf("unsupported OS")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download yt-dlp: unexpected status %s", resp.Status)
	}

	var destPath string
	if runtime.GOOS == "windows" {
		destPath = "bin/yt-dlp.exe"
	} else {
		destPath = "bin/yt-dlp"
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chmod(destPath, 0755)
	if err != nil {
		return err
	}

	binPath, err := filepath.Abs("bin")
	if err != nil {
		return err
	}
	path := os.Getenv("PATH")
	entries := strings.Split(path, string(os.PathListSeparator))
	for _, entry := range entries {
		if entry == binPath {
			return nil
		}
	}
	err = os.Setenv("PATH", binPath+string(os.PathListSeparator)+path)
	if err != nil {
		return err
	}

	return nil
}
