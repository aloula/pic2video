package renderjob

import (
"fmt"

"github.com/loula/pic2video/internal/domain/media"
)

func EvaluateQualityWarnings(job RenderJob, assets []media.Asset) []string {
	warnings := []string{}
	for _, a := range assets {
		if a.Width < job.Profile.Width || a.Height < job.Profile.Height {
			warnings = append(warnings, fmt.Sprintf("source image %s below target profile resolution", a.Path))
		}
	}
	return warnings
}
