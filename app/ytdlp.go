package app

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)



func binPaths() (ytdlp string, ffmpeg string) {
	if runtime.GOOS == "windows" {
		return "bin/yt-dlp.exe", "bin/ffmpeg.exe"
	}
	return "bin/yt-dlp", "bin/ffmpeg"
}




func runYtdlp(args []string) (stdout string, stderr string, err error) {
	os.MkdirAll("assets", 0755)

	ytdlp, ffmpeg := binPaths()

	logFile, err := os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", "", fmt.Errorf("creating log file: %w", err)
	}
	defer logFile.Close()

	if _, statErr := os.Stat(ffmpeg); statErr == nil {
		args = append(args, "--ffmpeg-location", ffmpeg)
	} else {
		fmt.Fprintf(logFile, "Warning: ffmpeg not found. Some features may not work correctly.\n")
	}

	var outBuf, errBuf strings.Builder
	cmd := exec.Command(ytdlp, args...)
	cmd.Stdout = io.MultiWriter(&outBuf, logFile)
	cmd.Stderr = io.MultiWriter(&errBuf, logFile)
	err = cmd.Run()

	return outBuf.String(), errBuf.String(), err
}


var throttleArgs = []string{
	"--sleep-requests", "1",
	"--sleep-interval", "5",
	"--max-sleep-interval", "10",
}
