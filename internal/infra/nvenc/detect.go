package nvenc

import (
"os/exec"
"strings"
)

func Available(ffmpegBin string) bool {
	if ffmpegBin == "" {
		ffmpegBin = "ffmpeg"
	}
	out, err := exec.Command(ffmpegBin, "-hide_banner", "-encoders").CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "h264_nvenc")
}
