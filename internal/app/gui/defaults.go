package gui

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ResolveLaunchDirectory() string {
	wd, err := os.Getwd()
	if err != nil || wd == "" {
		return "."
	}
	return wd
}

func DefaultConfiguration() GuiRunConfiguration {
	launchDir := ResolveLaunchDirectory()
	defaultOutputFolder := filepath.Join(launchDir, "output")
	return GuiRunConfiguration{
		InputFolder:     launchDir,
		OutputFolder:    defaultOutputFolder,
		Profile:         "uhd",
		ImageEffect:     "static",
		ImageDuration:   5,
		Transition:      1,
		FPS:             60,
		OrderMode:       "name",
		AudioSource:     "mp3",
		ExifFontSize:    42,
		Encoder:         "auto",
		Quality:         "high",
		Overwrite:       true,
		OutputFileName:  defaultOutputFilename("uhd", "high", 60),
		LaunchDirectory: launchDir,
	}
}

func defaultOutputFilename(profileName, quality string, fps int) string {
	profile := strings.ToLower(strings.TrimSpace(profileName))
	if profile != "fhd" {
		profile = "uhd"
	}
	q := strings.ToLower(strings.TrimSpace(quality))
	if q != "low" && q != "medium" {
		q = "high"
	}
	if fps <= 0 {
		fps = 60
	}
	return "slideshow_" + profile + "_" + q + "_" + strconv.Itoa(fps) + "fps.mp4"
}

func ResolveOutputPath(cfg GuiRunConfiguration) string {
	name := cfg.OutputFileName
	if name == "" {
		name = defaultOutputFilename(cfg.Profile, cfg.Quality, cfg.FPS)
	}
	if cfg.OutputFolder == "" {
		return filepath.Join(".", name)
	}
	return filepath.Join(cfg.OutputFolder, name)
}
