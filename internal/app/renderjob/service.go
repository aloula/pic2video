package renderjob

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/infra/ffmpeg"
	"github.com/loula/pic2video/internal/infra/fsio"
	"github.com/loula/pic2video/internal/infra/nvenc"
)

type Service struct{}

func FormatExifOverlayLine(exif *fsio.ExifData) string {
	if exif == nil {
		exif = &fsio.ExifData{}
	}
	return fmt.Sprintf(
		"%s - %s - %s - %s - %s - %s",
		fsio.NormalizeExifValue(exif.CameraModel),
		formatFocalDistance(exif.FocalDistance),
		formatSpeed(exif.ShutterSpeed),
		formatAperture(exif.Aperture),
		formatISO(exif.ISO),
		fsio.FormatCapturedDate(exif.CreateDate),
	)
}

func formatISO(v string) string {
	v = fsio.NormalizeExifValue(v)
	if strings.HasPrefix(strings.ToUpper(v), "ISO") {
		return v
	}
	return "ISO " + v
}

func parseRational(v string) (float64, bool) {
	v = strings.TrimSpace(v)
	parts := strings.Split(v, "/")
	if len(parts) != 2 {
		return 0, false
	}
	n, errN := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	d, errD := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if errN != nil || errD != nil || d == 0 {
		return 0, false
	}
	return n / d, true
}

func formatDecimal(v float64) string {
	if math.Abs(v-math.Round(v)) < 1e-9 {
		return strconv.FormatFloat(math.Round(v), 'f', 0, 64)
	}
	return strconv.FormatFloat(v, 'f', 1, 64)
}

func formatFocalDistance(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "Unknown"
	}
	if strings.HasSuffix(strings.ToLower(v), "mm") {
		return v
	}
	if fv, ok := parseRational(v); ok {
		return formatDecimal(fv) + "mm"
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return formatDecimal(mustFloat(v)) + "mm"
	}
	return v
}

func mustFloat(v string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
	return f
}

func formatSpeed(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "Unknown"
	}
	if strings.HasSuffix(v, "s") {
		return v
	}
	if fv, ok := parseRational(v); ok {
		if fv >= 1 {
			return formatDecimal(fv) + "s"
		}
		den := int(math.Round(1 / fv))
		if den > 0 {
			return fmt.Sprintf("1/%ds", den)
		}
	}
	if strings.HasPrefix(v, "1/") && strings.HasSuffix(v, "s") {
		return v
	}
	if strings.HasPrefix(v, "1/") {
		return v + "s"
	}
	return "1/" + v + "s"
}

func formatAperture(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "Unknown"
	}
	if strings.HasPrefix(strings.ToLower(v), "f/") {
		return v
	}
	if fv, ok := parseRational(v); ok {
		return "f/" + formatDecimal(fv)
	}
	return "f/" + v
}

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
	overlayLines := []string(nil)
	if job.ExifOverlayEnabled {
		overlayLines = make([]string, 0, len(job.InputAssets))
		for _, a := range job.InputAssets {
			exif, _ := fsio.ExtractExif(a.Path, job.FFprobeBin)
			overlayLines = append(overlayLines, FormatExifOverlayLine(exif))
		}
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio(
		job.OutputPath,
		job.InputAssets,
		job.AudioAssets,
		job.ImageEffect,
		job.ImageDurationSec,
		job.TransitionDurationSec,
		job.Profile.Width,
		job.Profile.Height,
		job.EffectiveEncoder,
		ffmpeg.OverlayOptions{
			Enabled:        job.ExifOverlayEnabled,
			FontSize:       job.ExifFontSize,
			FooterOffsetPx: job.ExifFooterOffsetPx,
			BoxAlpha:       job.ExifBoxAlpha,
			Lines:          overlayLines,
		},
	)
	if err := ffmpeg.Run(ctx, job.FFmpegBin, args); err != nil {
		ce := &ClassifiedError{Class: ErrExecution, Msg: "ffmpeg execution failed", Err: err}
		return BuildSummary(job, started, "failed", ce), ce
	}
	return BuildSummary(job, started, "success", nil), nil
}
