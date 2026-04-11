package e2e

import (
	"os"
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

func TestRenderMixedAspectExifOverlayMissingMetadataFallback(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "mixed-exif-overlay.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "uhd", "--exif-overlay", "--exif-font-size", "42", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	argsBytes, err := os.ReadFile(argsLog)
	if err != nil {
		t.Fatal(err)
	}
	args := string(argsBytes)
	if !strings.Contains(args, "drawtext=") {
		t.Fatalf("expected drawtext overlay args, got: %s", args)
	}
	if !strings.Contains(args, "Unknown - Unknown - Unknown - Unknown - Unknown -") {
		t.Fatalf("expected Unknown fallback values for missing exif metadata, got: %s", args)
	}
}
