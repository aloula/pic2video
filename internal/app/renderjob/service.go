package renderjob

import (
	"context"
	"os/exec"
	"time"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/infra/ffmpeg"
	"github.com/loula/pic2video/internal/infra/fsio"
	"github.com/loula/pic2video/internal/infra/nvenc"
)

type Service struct{}

func (s *Service) Run(ctx context.Context, job RenderJob) (RenderSummary, error) {
	started := time.Now()
	assets := pipeline.ApplyOrder(job.OrderMode, job.InputAssets, nil)
	if job.OrderMode == "explicit" && job.OrderFile != "" {
		explicit, err := fsio.ReadExplicitOrder(job.OrderFile)
		if err != nil {
			return BuildSummary(job, started, "failed", &ClassifiedError{Class: ErrInputValidation, Msg: "failed to read explicit order file", Err: err}), &ClassifiedError{Class: ErrInputValidation, Msg: "failed to read explicit order file", Err: err}
		}
		assets = pipeline.ApplyOrder(job.OrderMode, job.InputAssets, explicit)
	}
	job.InputAssets = assets
	for i, a := range job.InputAssets {
		p, err := ffmpeg.ProbeImage(job.FFprobeBin, a.Path)
		if err == nil {
			job.InputAssets[i].Width = p.Width
			job.InputAssets[i].Height = p.Height
		}
	}
	if err := ValidateJob(job); err != nil {
		return BuildSummary(job, started, "failed", err), err
	}
	if _, err := exec.LookPath(job.FFmpegBin); err == nil || job.FFmpegBin == "" {
		has := nvenc.Available(job.FFmpegBin)
		enc, err := nvenc.SelectEncoder(job.RequestedEncoder, has)
		if err != nil {
			ce := &ClassifiedError{Class: ErrEnvironment, Msg: "encoder selection failed", Err: err}
			return BuildSummary(job, started, "failed", ce), ce
		}
		job.EffectiveEncoder = enc
	} else {
		ce := &ClassifiedError{Class: ErrEnvironment, Msg: "ffmpeg binary unavailable", Err: err}
		return BuildSummary(job, started, "failed", ce), ce
	}
	job.Warnings = EvaluateQualityWarnings(job, job.InputAssets)
	args := ffmpeg.BuildRenderCommandArgsWithEffect(
		job.OutputPath,
		job.InputAssets,
		job.ImageEffect,
		job.ImageDurationSec,
		job.TransitionDurationSec,
		job.Profile.Width,
		job.Profile.Height,
		job.EffectiveEncoder,
	)
	if err := ffmpeg.Run(ctx, job.FFmpegBin, args); err != nil {
		ce := &ClassifiedError{Class: ErrExecution, Msg: "ffmpeg execution failed", Err: err}
		return BuildSummary(job, started, "failed", ce), ce
	}
	return BuildSummary(job, started, "success", nil), nil
}
