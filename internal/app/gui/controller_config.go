package gui

import (
	"strings"

	"math"
	"strconv"

	"fyne.io/fyne/v2/widget"
)

func CollectConfiguration(base GuiRunConfiguration, form *FormView, opts *OptionsView) GuiRunConfiguration {
	cfg := base
	cfg.InputFolder = strings.TrimSpace(form.InputEntry.Text)
	cfg.OutputFolder = strings.TrimSpace(form.OutputEntry.Text)
	cfg.Profile = strings.TrimSpace(opts.Profile.Selected)
	cfg.ImageEffect = strings.TrimSpace(opts.ImageEffect.Selected)
	if opts.ImageDur != nil {
		cfg.ImageDuration = math.Round(opts.ImageDur.Value)
	}
	if opts.Transition != nil {
		cfg.Transition = math.Round(opts.Transition.Value)
	}
	if opts.FPS30 != nil && opts.FPS30.Importance == widget.HighImportance {
		cfg.FPS = 30
	} else {
		cfg.FPS = 60
	}
	cfg.ExifOverlay = opts.ExifOverlay.Checked
	if v, err := strconv.Atoi(strings.TrimSpace(opts.ExifFontSize.Text)); err == nil {
		cfg.ExifFontSize = v
	}
	cfg.OrderMode = strings.TrimSpace(opts.OrderMode.Selected)
	if cfg.OrderMode == "explicit" {
		cfg.OrderFile = strings.TrimSpace(opts.OrderFile.Text)
	} else {
		cfg.OrderFile = ""
	}
	cfg.AudioSource = strings.TrimSpace(opts.AudioSource.Selected)
	if cfg.AudioSource == "" {
		cfg.AudioSource = "mp3"
	}
	switch {
	case opts.QualityLow != nil && opts.QualityLow.Importance == widget.HighImportance:
		cfg.Quality = "low"
	case opts.QualityMedium != nil && opts.QualityMedium.Importance == widget.HighImportance:
		cfg.Quality = "medium"
	default:
		cfg.Quality = "high"
	}
	cfg.Encoder = strings.TrimSpace(opts.Encoder.Selected)
	cfg.Overwrite = opts.Overwrite.Checked
	cfg.DebugExif = opts.DebugExif.Checked
	cfg.OutputFileName = defaultOutputFilename(cfg.Profile, cfg.Quality, cfg.FPS)
	cfg.OutputPath = ResolveOutputPath(cfg)
	return cfg
}
