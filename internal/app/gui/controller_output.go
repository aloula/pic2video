package gui

import (
	"fmt"
	"path/filepath"
	"strings"
)

func OutputPreviewText(cfg GuiRunConfiguration) string {
	return fmt.Sprintf("Output file: %s", compactPathTail(ResolveOutputPath(cfg), 3))
}

func compactPathTail(path string, parts int) string {
	cleaned := filepath.ToSlash(strings.TrimSpace(path))
	if cleaned == "" {
		return path
	}
	items := strings.Split(cleaned, "/")
	if len(items) <= parts {
		return cleaned
	}
	return ".../" + strings.Join(items[len(items)-parts:], "/")
}
