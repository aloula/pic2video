package renderjob

import (
"github.com/loula/pic2video/internal/domain/media"
"github.com/loula/pic2video/internal/domain/profile"
)

type BuildOptions struct {
	OutputPath      string
	ProfileName     string
	ImageDuration   float64
	Transition      float64
	Overwrite       bool
	OrderMode       string
	OrderFile       string
	RequestedEncode string
	FFmpegBin       string
	FFprobeBin      string
}

func BuildJob(opts BuildOptions, assets []media.Asset) (RenderJob, error) {
	p, err := profile.FromName(opts.ProfileName)
	if err != nil {
		return RenderJob{}, &ClassifiedError{Class: ErrInvalidArguments, Msg: "invalid profile", Err: err}
	}
	return RenderJob{
		InputAssets:           assets,
		OutputPath:            opts.OutputPath,
		Profile:               p,
		ImageDurationSec:      opts.ImageDuration,
		TransitionDurationSec: opts.Transition,
		Overwrite:             opts.Overwrite,
		OrderMode:             opts.OrderMode,
		OrderFile:             opts.OrderFile,
		RequestedEncoder:      opts.RequestedEncode,
		FFmpegBin:             opts.FFmpegBin,
		FFprobeBin:            opts.FFprobeBin,
	}, nil
}
