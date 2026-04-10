package e2e

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderExplicitOrderMode(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	manifest := filepath.Join(t.TempDir(), "order.txt")
	if err := os.WriteFile(manifest, []byte("c.jpg\na.jpg\nb.jpg\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "ordered.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--order", "explicit", "--order-file", manifest, "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}
