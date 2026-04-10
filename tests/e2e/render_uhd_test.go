package e2e

import (
	"path/filepath"
	"testing"
)

func TestRenderUHDHappyPath(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "uhd.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "uhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}
