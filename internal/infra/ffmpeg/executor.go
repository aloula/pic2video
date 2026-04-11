package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type executionError struct {
	err    error
	output string
}

func (e *executionError) Error() string {
	output := strings.TrimSpace(e.output)
	if output == "" {
		return e.err.Error()
	}
	return fmt.Sprintf("%v: %s", e.err, output)
}

func Run(ctx context.Context, ffmpegBin string, args []string) error {
	if ffmpegBin == "" {
		ffmpegBin = "ffmpeg"
	}
	cmd := exec.CommandContext(ctx, ffmpegBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &executionError{err: err, output: string(out)}
	}
	return err
}
