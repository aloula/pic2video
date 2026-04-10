package unit

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

var (
	unitCLIPath   string
	unitBuildOnce sync.Once
	unitBuildErr  error
)

func unitRepoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func unitCLIBinaryPath(t *testing.T) string {
	t.Helper()
	root := unitRepoRoot(t)
	unitBuildOnce.Do(func() {
		binDir, err := os.MkdirTemp("", "pic2video-unit-cli-")
		if err != nil {
			unitBuildErr = err
			return
		}
		unitCLIPath = filepath.Join(binDir, "pic2video")
		cmd := exec.Command("go", "build", "-o", unitCLIPath, "./cmd/pic2video")
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			unitBuildErr = &unitBuildError{err: err, output: string(out)}
		}
	})
	if unitBuildErr != nil {
		t.Fatal(unitBuildErr)
	}
	return unitCLIPath
}

func newUnitCLICommand(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(unitCLIBinaryPath(t), args...)
	cmd.Dir = unitRepoRoot(t)
	return cmd
}

type unitBuildError struct {
	err    error
	output string
}

func (e *unitBuildError) Error() string {
	return e.err.Error() + " output=" + e.output
}
