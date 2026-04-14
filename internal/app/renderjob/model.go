package renderjob

import (
	"time"

	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/domain/profile"
)

type RenderJob struct {
	InputAssets           []media.Asset
	AudioAssets           []string
	AudioSource           string
	OutputFPS             int
	ImageCount            int
	VideoCount            int
	OutputPath            string
	ExifOverlayEnabled    bool
	ExifFontSize          int
	ExifFooterOffsetPx    int
	ExifBoxAlpha          float64
	Profile               profile.Profile
	ImageEffect           string
	ImageDurationSec      float64
	TransitionDurationSec float64
	Overwrite             bool
	OrderMode             string
	OrderFile             string
	RequestedEncoder      string
	EffectiveEncoder      string
	Quality               string
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
	OutputFPS           int
	ImageCount          int
	VideoCount          int
	ExifOverlayEnabled  bool
	ExifFontSize        int
	ExifFooterOffsetPx  int
	ExifBoxAlpha        float64
	Status              string
	ErrorMessage        string
	Warnings            []string
}
