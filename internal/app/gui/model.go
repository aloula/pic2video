package gui

import "time"

type RunStatus string

const (
	RunStatusIdle         RunStatus = "Idle"
	RunStatusLoadingFiles RunStatus = "Loading files"
	RunStatusProcessing   RunStatus = "Processing:"
	RunStatusFinished     RunStatus = "Finished"
	RunStatusFailed       RunStatus = "Failed"
)

type GuiRunConfiguration struct {
	InputFolder     string
	OutputFolder    string
	Profile         string
	ImageEffect     string
	ImageDuration   float64
	Transition      float64
	FPS             int
	OrderMode       string
	OrderFile       string
	AudioSource     string
	ExifOverlay     bool
	ExifFontSize    int
	DebugExif       bool
	Encoder         string
	Quality         string
	Overwrite       bool
	FFmpegBin       string
	FFprobeBin      string
	OutputPath      string
	OutputFileName  string
	LaunchDirectory string
}

type GuiRunState struct {
	Status     RunStatus
	StartedAt  time.Time
	FinishedAt time.Time
	LastError  string
	ActivePID  int
}

type GuiLogEntry struct {
	Seq       int
	Timestamp time.Time
	Stream    string
	Message   string
}

type GuiValidationResult struct {
	OK                  bool
	Messages            []string
	SupportedMediaCount int
}
