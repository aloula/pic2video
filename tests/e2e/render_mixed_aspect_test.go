package e2e

import (
	"os"
	"strings"
	"testing"
)

func TestRenderMixedAspectQuality(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageAndVideoSet(t)
	cmd := newCLIRenderCommand(t, "--input", in, "--profile", "uhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	if !strings.Contains(string(outb), "warnings=") {
		t.Fatalf("expected warning summary in output: %s", string(outb))
	}
	if !strings.Contains(string(outb), "media: images=2 videos=1 fps=60") {
		t.Fatalf("expected mixed-media summary line with fps, got: %s", string(outb))
	}
}

func TestRenderMixedAspectWithSelectedFPS(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageAndVideoSet(t)
	cmd := newCLIRenderCommand(t, "--input", in, "--profile", "fhd", "--fps", "30", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	if !strings.Contains(string(outb), "media: images=2 videos=1 fps=30") {
		t.Fatalf("expected selected fps in summary, got: %s", string(outb))
	}
}

func TestRenderMixedAspectExifOverlayMissingMetadataFallback(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageSet(t)
	cmd := newCLIRenderCommand(t, "--input", in, "--profile", "uhd", "--exif-overlay", "--exif-font-size", "42", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
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
	if !strings.Contains(args, "Unknown - Unknown - Unknown - Unknown - ISO Unknown -") {
		t.Fatalf("expected Unknown fallback values for missing exif metadata, got: %s", args)
	}
}
