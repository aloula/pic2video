package renderjob

import (
	"time"

	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/domain/profile"
)

type RenderJob struct {
	InputAssets           []media.Asset
	OutputPath            string
	Profile               profile.Profile
	ImageEffect           string
	ImageDurationSec      float64
	TransitionDurationSec float64
	Overwrite             bool
	OrderMode             string
	OrderFile             string
	RequestedEncoder      string
	EffectiveEncoder      string
	Warnings              []string
	FFmpegBin             string
	FFprobeBin            string
}

type RenderSummary struct {
	JobID               string
	StartedAt           time.Time
	FinishedAt          time.Time
	ElapsedSeconds      float64
	ProcessedAssets     int
	SkippedAssets       int
	ProfileName         string
	EffectiveResolution string
	EffectiveEncoder    string
	OutputPath          string
	Status              string
	ErrorMessage        string
	Warnings            []string
}
