package renderjob

import (
	"fmt"
	"time"
)

func BuildSummary(job RenderJob, started time.Time, status string, err error) RenderSummary {
	finished := time.Now()
	r := RenderSummary{
		JobID:               started.Format("20060102-150405"),
		StartedAt:           started,
		FinishedAt:          finished,
		ElapsedSeconds:      finished.Sub(started).Seconds(),
		ProcessedAssets:     len(job.InputAssets),
		ProfileName:         job.Profile.Name,
		EffectiveResolution: fmt.Sprintf("%dx%d", job.Profile.Width, job.Profile.Height),
		EffectiveEncoder:    job.EffectiveEncoder,
		OutputPath:          job.OutputPath,
		ExifOverlayEnabled:  job.ExifOverlayEnabled,
		ExifFontSize:        job.ExifFontSize,
		ExifFooterOffsetPx:  job.ExifFooterOffsetPx,
		ExifBoxAlpha:        job.ExifBoxAlpha,
		Warnings:            job.Warnings,
	}
	if err != nil {
		r.Status = "failed"
		r.ErrorMessage = err.Error()
	} else {
		r.Status = status
	}
	return r
}
