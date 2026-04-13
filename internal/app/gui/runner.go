package gui

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type OutputHandler func(stream, line string)

type Runner struct {
	CLIBin string
}

func NewRunner(cliBin string) *Runner {
	cliBin = resolveCLIBin(cliBin)
	return &Runner{CLIBin: cliBin}
}

func resolveCLIBin(cliBin string) string {
	if strings.TrimSpace(cliBin) != "" {
		return cliBin
	}
	if env := strings.TrimSpace(os.Getenv("P2V_CLI_BIN")); env != "" {
		return env
	}
	if resolved, ok := findNearbyCLIBinary(); ok {
		return resolved
	}
	if lp, err := exec.LookPath("pic2video"); err == nil {
		return lp
	}
	return "pic2video"
}

func findNearbyCLIBinary() (string, bool) {
	exePath, err := os.Executable()
	if err != nil || strings.TrimSpace(exePath) == "" {
		return "", false
	}
	exeDir := filepath.Dir(exePath)
	dirs := []string{
		exeDir,
		filepath.Clean(filepath.Join(exeDir, "..", "cli")),
		filepath.Clean(filepath.Join(exeDir, "..")),
	}
	// Cover local host and cross-build naming conventions used in Makefile outputs.
	names := []string{
		"pic2video",
		"pic2video.exe",
		"pic2video-linux-amd64",
		"pic2video-darwin-amd64",
		"pic2video-windows-amd64.exe",
	}
	for _, d := range dirs {
		for _, n := range names {
			candidate := filepath.Clean(filepath.Join(d, n))
			if st, statErr := os.Stat(candidate); statErr == nil && !st.IsDir() {
				return candidate, true
			}
		}
	}
	return "", false
}

func (r *Runner) Run(ctx context.Context, cfg GuiRunConfiguration, onState func(RunStatus), onOutput OutputHandler) error {
	if onState != nil {
		onState(RunStatusLoadingFiles)
	}
	args := BuildRenderCommandArgs(cfg)
	cmd := exec.CommandContext(ctx, r.CLIBin, args...)
	if cfg.LaunchDirectory != "" {
		cmd.Dir = cfg.LaunchDirectory
	} else {
		cmd.Dir = "."
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if onOutput != nil {
		onOutput("system", fmt.Sprintf("launch: %s %v", filepath.Clean(r.CLIBin), args))
	}
	if err := cmd.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) || strings.Contains(strings.ToLower(err.Error()), "executable file not found") {
			return fmt.Errorf("failed to start render command: %w (resolved cli binary=%q, set P2V_CLI_BIN to explicit path if needed)", err, r.CLIBin)
		}
		return fmt.Errorf("failed to start render command: %w", err)
	}
	if onState != nil {
		onState(RunStatusProcessing)
	}

	var wg sync.WaitGroup
	read := func(stream string, rc io.Reader) {
		defer wg.Done()
		s := bufio.NewScanner(rc)
		for s.Scan() {
			if onOutput != nil {
				onOutput(stream, s.Text())
			}
		}
	}
	wg.Add(2)
	go read("stdout", stdout)
	go read("stderr", stderr)

	waitErr := cmd.Wait()
	wg.Wait()
	if waitErr != nil {
		if onState != nil {
			onState(RunStatusFailed)
		}
		return waitErr
	}
	if onState != nil {
		onState(RunStatusFinished)
	}
	return nil
}
