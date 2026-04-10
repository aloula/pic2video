package e2e

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderMixedAspectQuality(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "mixed.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "uhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	if !strings.Contains(string(outb), "warnings=") {
		t.Fatalf("expected warning summary in output: %s", string(outb))
	}
}
