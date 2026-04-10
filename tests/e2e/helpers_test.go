package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

var (
	builtCLIPath string
	buildCLIOnce sync.Once
	buildCLIErr  error
)

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func cliBinaryPath(t *testing.T) string {
	t.Helper()
	root := repoRoot(t)
	buildCLIOnce.Do(func() {
		binDir, err := os.MkdirTemp("", "pic2video-cli-")
		if err != nil {
			buildCLIErr = err
			return
		}
		builtCLIPath = filepath.Join(binDir, "pic2video")
		cmd := exec.Command("go", "build", "-o", builtCLIPath, "./cmd/pic2video")
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			buildCLIErr = &buildError{err: err, output: string(out)}
		}
	})
	if buildCLIErr != nil {
		t.Fatal(buildCLIErr)
	}
	return builtCLIPath
}

func newCLIRenderCommand(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	all := append([]string{"render"}, args...)
	cmd := exec.Command(cliBinaryPath(t), all...)
	cmd.Dir = repoRoot(t)
	return cmd
}

type buildError struct {
	err    error
	output string
}

func (e *buildError) Error() string {
	return e.err.Error() + " output=" + e.output
}

func createFakeBinaries(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	ffmpeg := filepath.Join(dir, "ffmpeg")
	ffprobe := filepath.Join(dir, "ffprobe")
	if err := os.WriteFile(ffmpeg, []byte("#!/bin/sh\nif [ \"$2\" = \"-encoders\" ] || [ \"$1\" = \"-hide_banner\" ]; then\n  echo ' V..... h264_nvenc NVIDIA NVENC H.264 encoder'\n  exit 0\nfi\nfor last; do :; done\nout=\"$last\"\nmkdir -p \"$(dirname \"$out\")\"\ntouch \"$out\"\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ffprobe, []byte("#!/bin/sh\necho '{\"streams\":[{\"width\":4000,\"height\":3000,\"codec_name\":\"mjpeg\"}]}'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return ffmpeg, ffprobe
}

func createImageSet(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, n := range []string{"a.jpg", "b.jpg", "c.jpg"} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}
