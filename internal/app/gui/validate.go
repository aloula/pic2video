package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/loula/pic2video/internal/infra/fsio"
)

func ValidatePreflight(cfg GuiRunConfiguration) GuiValidationResult {
	res := GuiValidationResult{OK: true, Messages: make([]string, 0, 4)}

	if strings.TrimSpace(cfg.InputFolder) == "" {
		res.OK = false
		res.Messages = append(res.Messages, "Input folder is required")
		return res
	}
	assets, err := fsio.ListMixedAssets(cfg.InputFolder)
	if err != nil {
		res.OK = false
		res.Messages = append(res.Messages, fmt.Sprintf("Input folder validation failed: %v", err))
	} else {
		res.SupportedMediaCount = len(assets)
		if len(assets) == 0 {
			res.OK = false
			res.Messages = append(res.Messages, "Input folder has no supported media files")
		}
	}

	if strings.TrimSpace(cfg.OutputFolder) == "" {
		res.OK = false
		res.Messages = append(res.Messages, "Output folder is required")
		return res
	}
	if err := os.MkdirAll(cfg.OutputFolder, 0o755); err != nil {
		res.OK = false
		res.Messages = append(res.Messages, fmt.Sprintf("Output folder is not writable: %v", err))
		return res
	}
	tmp, err := os.CreateTemp(cfg.OutputFolder, "p2v-gui-writecheck-*.tmp")
	if err != nil {
		res.OK = false
		res.Messages = append(res.Messages, fmt.Sprintf("Output folder is not writable: %v", err))
		return res
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	_ = os.Remove(tmpPath)

	// Normalize for downstream command construction.
	cfg.OutputFolder = filepath.Clean(cfg.OutputFolder)
	return res
}
