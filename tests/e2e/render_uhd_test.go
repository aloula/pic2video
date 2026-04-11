package e2e

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRenderUHDKenBurnsMedium(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "uhd-kenburns-medium.mp4")
	cmd := newCLIRenderCommand(t,
		"--input", in,
		"--output", out,
		"--profile", "uhd",
		"--image-effect", "kenburns-medium",
		"--ffmpeg-bin", ffmpeg,
		"--ffprobe-bin", ffprobe,
	)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}

func TestRenderUHDWithExifOverlayArgs(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "uhd-exif-overlay.mp4")
	cmd := newCLIRenderCommand(t,
		"--input", in,
		"--output", out,
		"--profile", "uhd",
		"--exif-overlay",
		"--exif-font-size", "60",
		"--ffmpeg-bin", ffmpeg,
		"--ffprobe-bin", ffprobe,
	)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	argsBytes, err := os.ReadFile(argsLog)
	if err != nil {
		t.Fatal(err)
	}
	args := string(argsBytes)
	if !strings.Contains(args, "drawtext=") {
		t.Fatalf("expected drawtext filter with overlay enabled, got: %s", args)
	}
	if !strings.Contains(args, "fontsize=60") {
		t.Fatalf("expected requested font size in drawtext filter, got: %s", args)
	}
	if !strings.Contains(args, "y=h-th-60") {
		t.Fatalf("expected footer offset in drawtext filter, got: %s", args)
	}
	if !strings.Contains(args, "fontcolor=white") || !strings.Contains(args, "boxcolor=black@0.40") {
		t.Fatalf("expected overlay text color and background transparency style, got: %s", args)
	}
}
