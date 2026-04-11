package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/loula/pic2video/internal/infra/nvenc"
)

// StartOptions captures user-visible render options for startup status output.
type StartOptions struct {
	Input              string
	Output             string
	Profile            string
	ImageEffect        string
	ImageDuration      float64
	TransitionDuration float64
	Order              string
	OrderFile          string
	AudioFiles         int
	AudioOrder         string
	ExifOverlay        bool
	ExifFontSize       int
	ExifFooterOffsetPx int
	ExifBoxAlpha       float64
	Encoder            string
	Overwrite          bool
	Files              int
}

func FormatOutputFormat(outputPath string) string {
	ext := strings.TrimPrefix(filepath.Ext(outputPath), ".")
	if ext == "" {
		return "UNKNOWN"
	}
	return strings.ToUpper(ext)
}

func FormatElapsed(seconds float64) string {
	if seconds < 1.0 {
		return "< 1s"
	}
	if seconds < 60.0 {
		return fmt.Sprintf("%.1fs", seconds)
	}
	totalSeconds := int(seconds)
	minutes := totalSeconds / 60
	remaining := totalSeconds % 60
	return fmt.Sprintf("%dm %ds", minutes, remaining)
}

func FormatAnnouncement(opts StartOptions) string {
	orderFile := opts.OrderFile
	if strings.TrimSpace(orderFile) == "" {
		orderFile = "-"
	}
	audioOrder := opts.AudioOrder
	if strings.TrimSpace(audioOrder) == "" {
		audioOrder = "-"
	}

	return fmt.Sprintf(
		"status=starting files=%d format=%s\n"+
			"details: input=%s output=%s profile=%s effect=%s encoder=%s overwrite=%t\n"+
			"timing: image-duration=%.1fs transition-duration=%.1fs\n"+
			"order: mode=%s order-file=%s\n"+
			"audio: files=%d order=%s\n"+
			"exif-overlay: enabled=%t font-size=%d footer-offset=%d box-alpha=%.2f",
		opts.Files,
		FormatOutputFormat(opts.Output),
		opts.Input,
		opts.Output,
		opts.Profile,
		opts.ImageEffect,
		opts.Encoder,
		opts.Overwrite,
		opts.ImageDuration,
		opts.TransitionDuration,
		opts.Order,
		orderFile,
		opts.AudioFiles,
		audioOrder,
		opts.ExifOverlay,
		opts.ExifFontSize,
		opts.ExifFooterOffsetPx,
		opts.ExifBoxAlpha,
	)
}

func FormatSummary(profile, res string, exifOverlay bool, exifFontSize, exifFooterOffsetPx int, exifBoxAlpha float64, requested, effective, output string, elapsed float64, processed int, warnings []string, nvencAvailable bool) string {
	report := nvenc.BuildReport(requested, effective, nvencAvailable)
	return fmt.Sprintf(
		"status=success\n"+
			"result: profile=%s resolution=%s %s processed=%d files=%d\n"+
			"exif-overlay: enabled=%t font-size=%d footer-offset=%d box-alpha=%.2f\n"+
			"output: format=%s elapsed=%s output=%s warnings=%d",
		profile,
		res,
		report,
		processed,
		processed,
		exifOverlay,
		exifFontSize,
		exifFooterOffsetPx,
		exifBoxAlpha,
		FormatOutputFormat(output),
		FormatElapsed(elapsed),
		output,
		len(warnings),
	)
}
