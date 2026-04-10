package e2e

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSmoke(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "smoke.mp4")

	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render failed: %v output=%s", err, string(outb))
	}
	output := string(outb)

	if !strings.Contains(output, "status=starting") {
		t.Fatalf("expected pre-render announcement: %s", output)
	}
	if !strings.Contains(output, "details:") {
		t.Fatalf("expected details section in startup announcement: %s", output)
	}
	if !strings.Contains(output, "timing:") {
		t.Fatalf("expected timing section in startup announcement: %s", output)
	}
	if !strings.Contains(output, "order:") {
		t.Fatalf("expected order section in startup announcement: %s", output)
	}
	if !strings.Contains(output, "input=") || !strings.Contains(output, "output=") || !strings.Contains(output, "profile=") || !strings.Contains(output, "encoder=") || !strings.Contains(output, "overwrite=") {
		t.Fatalf("expected full startup options in details section: %s", output)
	}
	if !strings.Contains(output, "image-duration=") || !strings.Contains(output, "transition-duration=") {
		t.Fatalf("expected timing options in startup announcement: %s", output)
	}
	if !strings.Contains(output, "mode=") || !strings.Contains(output, "order-file=") {
		t.Fatalf("expected order options in startup announcement: %s", output)
	}
	if !strings.Contains(output, "status=success") {
		t.Fatalf("expected completion status: %s", output)
	}
	if !strings.Contains(output, "files=") {
		t.Fatalf("expected files field in output: %s", output)
	}
	if !strings.Contains(output, "format=MP4") {
		t.Fatalf("expected MP4 format field in output: %s", output)
	}
	if !strings.Contains(output, "processed=") {
		t.Fatalf("expected processed field for backward compatibility: %s", output)
	}
	if !strings.Contains(output, "elapsed=") {
		t.Fatalf("expected elapsed field in output: %s", output)
	}

	startingIdx := strings.Index(output, "status=starting")
	successIdx := strings.Index(output, "status=success")
	if startingIdx == -1 || successIdx == -1 || startingIdx > successIdx {
		t.Fatalf("expected status=starting before status=success: %s", output)
	}
}

func TestSmokeAnnouncementTiming(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	in := createImageSet(t)
	out := filepath.Join(t.TempDir(), "timing.mp4")

	cmd := newCLIRenderCommand(t, "--input", in, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		_ = cmd.Wait()
		t.Fatal("expected first stdout line")
	}
	firstLine := scanner.Text()
	delta := time.Since(start)
	if !strings.Contains(firstLine, "status=starting") {
		_ = cmd.Wait()
		t.Fatalf("expected first line to be status=starting, got: %s", firstLine)
	}
	if delta > time.Second {
		_ = cmd.Wait()
		t.Fatalf("expected status=starting within 1s, got %s", delta)
	}

	for scanner.Scan() {
	}
	if err := scanner.Err(); err != nil {
		_ = cmd.Wait()
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("render failed: %v", err)
	}
}

func TestSmokeNoCompletionOnFailure(t *testing.T) {
	ffmpeg, ffprobe := createFakeBinaries(t)
	out := filepath.Join(t.TempDir(), "failed.mp4")
	missingInput := filepath.Join(t.TempDir(), "missing")

	cmd := newCLIRenderCommand(t, "--input", missingInput, "--output", out, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	outb, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit")
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 0 {
		t.Fatal("expected non-zero exit code")
	}
	output := string(outb)
	if strings.Contains(output, "status=success") {
		t.Fatalf("did not expect success status on failure: %s", output)
	}
	if strings.Contains(output, "files=") || strings.Contains(output, "format=") || strings.Contains(output, "elapsed=") {
		t.Fatalf("did not expect completion fields on failure: %s", output)
	}
}
