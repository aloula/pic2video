package ffmpeg

import (
"context"
"os/exec"
)

func Run(ctx context.Context, ffmpegBin string, args []string) error {
	if ffmpegBin == "" {
		ffmpegBin = "ffmpeg"
	}
	cmd := exec.CommandContext(ctx, ffmpegBin, args...)
	_, err := cmd.CombinedOutput()
	return err
}
