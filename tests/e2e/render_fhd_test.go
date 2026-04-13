package e2e

import (
	"os"
	"strings"
	"testing"
)

func TestRenderFHDHappyPath(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	cmd := newCLIRenderCommand(t, "--input", in, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}

func TestRenderFHDKenBurnsMedium(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	cmd := newCLIRenderCommand(t,
		"--input", in,
		"--profile", "fhd",
		"--image-effect", "kenburns-medium",
		"--ffmpeg-bin", ffmpeg,
		"--ffprobe-bin", ffprobe,
	)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
}

func TestRenderFHDWithMP3AudioOrdered(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageAndAudioSet(t)
	cmd := newCLIRenderCommand(t, "--input", in, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	argsBytes, err := os.ReadFile(argsLog)
	if err != nil {
		t.Fatal(err)
	}
	args := string(argsBytes)
	aIdx := strings.Index(args, "ambient_a.mp3")
	bIdx := strings.Index(args, "ambient_b.mp3")
	if aIdx == -1 || bIdx == -1 || aIdx > bIdx {
		t.Fatalf("expected alphabetical mp3 order in ffmpeg args, got: %s", args)
	}
	if strings.Contains(args, "ignored.wav") {
		t.Fatalf("expected non-mp3 audio to be ignored, got: %s", args)
	}
	if strings.Contains(args, "-stream_loop\n-1") {
		t.Fatalf("did not expect per-input mp3 looping; next mp3 should play after the first, got: %s", args)
	}
	if !strings.Contains(args, "[aout]") {
		t.Fatalf("expected mapped audio output in ffmpeg args, got: %s", args)
	}
}

func TestRenderFHDFadeEnabledOutputGeneration(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageSet(t)
	cmd := newCLIRenderCommand(t, "--input", in, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	argsBytes, err := os.ReadFile(argsLog)
	if err != nil {
		t.Fatal(err)
	}
	args := string(argsBytes)
	if !strings.Contains(args, "fade=t=in:st=0") || !strings.Contains(args, "fade=t=out:st=") {
		t.Fatalf("expected fade-in and fade-out directives in ffmpeg args, got: %s", args)
	}
}

func TestRenderFHDMixedMediaArgsStableAcrossRuns(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageAndAudioSet(t)
	run := func() string {
		cmd := newCLIRenderCommand(t, "--input", in, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
		if outb, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("render failed: %v output=%s", err, string(outb))
		}
		argsBytes, err := os.ReadFile(argsLog)
		if err != nil {
			t.Fatal(err)
		}
		return string(argsBytes)
	}
	first := run()
	second := run()
	if first != second {
		t.Fatalf("expected deterministic ffmpeg args across runs\nfirst=%s\nsecond=%s", first, second)
	}
}

func TestRenderFHDWithExifOverlayArgs(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageSet(t)
	cmd := newCLIRenderCommand(t,
		"--input", in,
		"--profile", "fhd",
		"--exif-overlay",
		"--exif-font-size", "42",
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
		t.Fatalf("expected drawtext filter with exif overlay enabled, got: %s", args)
	}
	if !strings.Contains(args, "fontsize=42") {
		t.Fatalf("expected requested font size in drawtext filter, got: %s", args)
	}
	if !strings.Contains(args, "y=h-th-30") {
		t.Fatalf("expected footer offset in drawtext filter, got: %s", args)
	}
	if !strings.Contains(args, "fontcolor=white") || !strings.Contains(args, "boxcolor=black@0.40") {
		t.Fatalf("expected text color and transparency style in drawtext filter, got: %s", args)
	}
	if !strings.Contains(args, "Unknown - Unknown - Unknown - Unknown - ISO Unknown -") {
		t.Fatalf("expected deterministic EXIF overlay field order with Unknown fallback, got: %s", args)
	}
}

func TestRenderFHDWithExifOverlayFontBoundaries(t *testing.T) {
	for _, size := range []string{"36", "60"} {
		t.Run(size, func(t *testing.T) {
			ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
			in := createImageSet(t)
			cmd := newCLIRenderCommand(t,
				"--input", in,
				"--profile", "fhd",
				"--exif-overlay",
				"--exif-font-size", size,
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
			if !strings.Contains(string(argsBytes), "fontsize="+size) {
				t.Fatalf("expected fontsize=%s in ffmpeg args, got: %s", size, string(argsBytes))
			}
		})
	}
}

func TestRenderFHDWithSelectedFPS(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageAndVideoSet(t)
	cmd := newCLIRenderCommand(t,
		"--input", in,
		"--profile", "fhd",
		"--fps", "30",
		"--ffmpeg-bin", ffmpeg,
		"--ffprobe-bin", ffprobe,
	)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	if !strings.Contains(string(outb), "media: images=2 videos=1 fps=30") {
		t.Fatalf("expected selected fps summary output, got: %s", string(outb))
	}
}
