package gui

import (
	"strconv"
	"strings"
)

func BuildRenderCommandArgs(cfg GuiRunConfiguration) []string {
	args := []string{"render"}

	in := strings.TrimSpace(cfg.InputFolder)
	if in == "" {
		in = "."
	}
	args = append(args, "--input", in)

	profile := strings.TrimSpace(strings.ToLower(cfg.Profile))
	if profile == "" {
		profile = "uhd"
	}
	args = append(args, "--profile", profile)

	if cfg.ImageEffect != "" {
		args = append(args, "--image-effect", cfg.ImageEffect)
	}
	if cfg.ImageDuration > 0 {
		args = append(args, "--image-duration", strconv.FormatFloat(cfg.ImageDuration, 'f', -1, 64))
	}
	if cfg.Transition >= 0 {
		args = append(args, "--transition-duration", strconv.FormatFloat(cfg.Transition, 'f', -1, 64))
	}
	if cfg.FPS > 0 {
		args = append(args, "--fps", strconv.Itoa(cfg.FPS))
	}

	order := strings.TrimSpace(cfg.OrderMode)
	if order == "" {
		order = "name"
	}
	args = append(args, "--order", order)
	audioSource := strings.TrimSpace(cfg.AudioSource)
	if audioSource == "" {
		audioSource = "mp3"
	}
	args = append(args, "--audio-source", audioSource)
	if strings.TrimSpace(cfg.OrderFile) != "" {
		args = append(args, "--order-file", cfg.OrderFile)
	}

	if cfg.ExifOverlay {
		args = append(args, "--exif-overlay")
	}
	if cfg.ExifFontSize > 0 {
		args = append(args, "--exif-font-size", strconv.Itoa(cfg.ExifFontSize))
	}
	if cfg.DebugExif {
		args = append(args, "--debug-exif")
	}

	enc := strings.TrimSpace(cfg.Encoder)
	if enc == "" {
		enc = "auto"
	}
	args = append(args, "--encoder", enc)
	args = append(args, "--overwrite", strconv.FormatBool(cfg.Overwrite))

	if strings.TrimSpace(cfg.FFmpegBin) != "" {
		args = append(args, "--ffmpeg-bin", cfg.FFmpegBin)
	}
	if strings.TrimSpace(cfg.FFprobeBin) != "" {
		args = append(args, "--ffprobe-bin", cfg.FFprobeBin)
	}

	return args
}
