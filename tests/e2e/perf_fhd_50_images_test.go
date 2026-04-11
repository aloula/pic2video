package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPerfFHD50Images(t *testing.T) {
	if os.Getenv("RUN_PERF") == "" {
		t.Skip("set RUN_PERF=1 to run performance benchmark assertion")
	}
	ffmpeg, ffprobe := createFakeBinaries(t)
	dir := t.TempDir()
	for i := 0; i < 50; i++ {
		name := filepath.Join(dir, fmt.Sprintf("img-%03d.jpg", i))
		if err := os.WriteFile(name, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	start := time.Now()
	cmd := newCLIRenderCommand(t, "--input", dir, "--profile", "fhd", "--ffmpeg-bin", ffmpeg, "--ffprobe-bin", ffprobe)
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("perf render failed: %v output=%s", err, string(outb))
	}
	if time.Since(start) > 5*time.Minute {
		t.Fatalf("benchmark exceeded 5 minutes")
	}
}
