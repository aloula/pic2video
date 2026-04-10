package e2e

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMissingFFmpegExitCode(t *testing.T) {
	_, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "out.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", filepath.Join(t.TempDir(), "not-found"), "--ffprobe-bin", ffprobe)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() != 4 {
		t.Fatalf("expected exit code 4, got %d", ee.ExitCode())
	}
}
