package renderjob

import (
	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/domain/profile"
)

type BuildOptions struct {
	OutputPath         string
	AudioAssets        []string
	AudioSource        string
	OutputFPS          int
	ExifOverlay        bool
	ExifFontSize       int
	ExifFooterOffsetPx int
	ExifBoxAlpha       float64
	ProfileName        string
	ImageEffect        string
	ImageDuration      float64
	Transition         float64
	Overwrite          bool
	OrderMode          string
	OrderFile          string
	RequestedEncode    string
	FFmpegBin          string
	FFprobeBin         string
	Quality            string
}

func BuildJob(opts BuildOptions, assets []media.Asset) (RenderJob, error) {
	p, err := profile.FromName(opts.ProfileName)
	if err != nil {
		return RenderJob{}, &ClassifiedError{Class: ErrInvalidArguments, Msg: "invalid profile", Err: err}
	}
	fps := opts.OutputFPS
	if fps <= 0 {
		fps = 60
	}
	imageCount := 0
	videoCount := 0
	audioSource := opts.AudioSource
	if audioSource == "" {
		audioSource = "mp3"
	}
	quality := opts.Quality
	if quality == "" {
		quality = "high"
	}
	for _, a := range assets {
		if a.MediaType == media.MediaTypeVideo {
			videoCount++
		} else {
			imageCount++
		}
	}
	return RenderJob{
		InputAssets:           assets,
		AudioAssets:           opts.AudioAssets,
		AudioSource:           audioSource,
		OutputFPS:             fps,
		ImageCount:            imageCount,
		VideoCount:            videoCount,
		OutputPath:            opts.OutputPath,
		ExifOverlayEnabled:    opts.ExifOverlay,
		ExifFontSize:          opts.ExifFontSize,
		ExifFooterOffsetPx:    opts.ExifFooterOffsetPx,
		ExifBoxAlpha:          opts.ExifBoxAlpha,
		Profile:               p,
		ImageEffect:           opts.ImageEffect,
		ImageDurationSec:      opts.ImageDuration,
		TransitionDurationSec: opts.Transition,
		Overwrite:             opts.Overwrite,
		OrderMode:             opts.OrderMode,
		OrderFile:             opts.OrderFile,
		RequestedEncoder:      opts.RequestedEncode,
		FFmpegBin:             opts.FFmpegBin,
		FFprobeBin:            opts.FFprobeBin,
		Quality:               quality,
	}, nil
}
