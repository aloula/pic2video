package e2e

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInvalidInputExitCode(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	out := filepath.Join(t.TempDir(), "bad.mp4")
	cmd := newCLIRenderCommand(t, "--input", filepath.Join(t.TempDir(), "missing"), "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() != 3 {
		t.Fatalf("expected exit code 3, got %d", ee.ExitCode())
	}
}

func TestInvalidExifFontSizeExitCode(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "bad-font.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--exif-overlay", "--exif-font-size", "20", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() != 2 {
		t.Fatalf("expected exit code 2, got %d", ee.ExitCode())
	}
}
