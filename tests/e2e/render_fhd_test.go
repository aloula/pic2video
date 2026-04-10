package e2e

import (
	"path/filepath"
	"testing"
)

func TestRenderFHDHappyPath(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "fhd.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}

func TestRenderFHDKenBurnsMedium(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "fhd-kenburns-medium.mp4")
	cmd := newCLIRenderCommand(t,
		"--input", in,
		"--output", out,
		"--profile", "fhd",
		"--image-effect", "kenburns-medium",
		"--ffmpeg-bin", ffmpeg,
		"--ffprobe-bin", ffprobe,
	)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}
