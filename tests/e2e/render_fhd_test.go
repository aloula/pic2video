package e2e

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRenderFHDWithMP3AudioOrdered(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageAndAudioSet(t)
	out := filepath.Join(t.TempDir(), "fhd-audio.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
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
	if !strings.Contains(args, "[aout]") {
		t.Fatalf("expected mapped audio output in ffmpeg args, got: %s", args)
	}
}

func TestRenderFHDFadeEnabledOutputGeneration(t *testing.T) {
	ffmpeg, ffprobe, argsLog := createFakeBinariesWithArgsCapture(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "fhd-fades.mp4")
	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
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
	out := filepath.Join(t.TempDir(), "fhd-stable.mp4")
	run := func() string {
		cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
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
