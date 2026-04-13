package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/loula/pic2video/internal/app/gui"
)

func writeFakeGUIRunnerBinary(t *testing.T, script string) string {
	t.Helper()
	d := t.TempDir()
	p := filepath.Join(d, "pic2video")
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGUIStatusFlowSuccess(t *testing.T) {
	bin := writeFakeGUIRunnerBinary(t, "#!/bin/sh\necho status=starting\necho status=processing\nexit 0\n")
	r := gui.NewRunner(bin)
	cfg := gui.DefaultConfiguration()
	cfg.InputFolder = createImageSet(t)
	cfg.OutputFolder = t.TempDir()
	seen := make([]gui.RunStatus, 0, 4)
	err := r.Run(context.Background(), cfg, func(s gui.RunStatus) { seen = append(seen, s) }, func(_, _ string) {})
	if err != nil {
		t.Fatalf("expected successful runner execution, got %v", err)
	}
	if len(seen) < 2 || seen[0] != gui.RunStatusLoadingFiles || seen[len(seen)-1] != gui.RunStatusFinished {
		t.Fatalf("unexpected status flow: %+v", seen)
	}
}

func TestGUIStatusFlowFailure(t *testing.T) {
	bin := writeFakeGUIRunnerBinary(t, "#!/bin/sh\necho status=starting\necho error=boom 1>&2\nexit 1\n")
	r := gui.NewRunner(bin)
	cfg := gui.DefaultConfiguration()
	cfg.InputFolder = createImageSet(t)
	cfg.OutputFolder = t.TempDir()
	seen := make([]gui.RunStatus, 0, 4)
	err := r.Run(context.Background(), cfg, func(s gui.RunStatus) { seen = append(seen, s) }, func(_, _ string) {})
	if err == nil {
		t.Fatal("expected runner failure")
	}
	if len(seen) == 0 || seen[len(seen)-1] != gui.RunStatusFailed {
		t.Fatalf("expected failed terminal status, got %+v", seen)
	}
}
