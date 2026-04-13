package e2e

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/gui"
)

func TestGUILogFlowCapturesRunnerOutput(t *testing.T) {
	d := t.TempDir()
	bin := filepath.Join(d, "pic2video")
	script := "#!/bin/sh\necho first-line\necho second-line 1>&2\nexit 0\n"
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	r := gui.NewRunner(bin)
	cfg := gui.DefaultConfiguration()
	cfg.InputFolder = createImageSet(t)
	cfg.OutputFolder = t.TempDir()
	store := gui.NewLogStore(20)
	err := r.Run(context.Background(), cfg, nil, store.Append)
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	txt := store.Text()
	if !strings.Contains(txt, "first-line") || !strings.Contains(txt, "second-line") {
		t.Fatalf("expected captured logs, got: %s", txt)
	}
}
