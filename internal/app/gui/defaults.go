package gui

import (
	"os"
	"path/filepath"
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
		Overwrite:       true,
		OutputFileName:  defaultOutputFilename("uhd"),
		LaunchDirectory: launchDir,
	}
}

func defaultOutputFilename(profileName string) string {
	if profileName == "fhd" {
		return "slideshow_fhd.mp4"
	}
	return "slideshow_uhd.mp4"
}

func ResolveOutputPath(cfg GuiRunConfiguration) string {
	name := cfg.OutputFileName
	if name == "" {
		name = defaultOutputFilename(cfg.Profile)
	}
	if cfg.OutputFolder == "" {
		return filepath.Join(".", name)
	}
	return filepath.Join(cfg.OutputFolder, name)
}
