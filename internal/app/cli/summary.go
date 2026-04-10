package cli

import (
"fmt"

"github.com/loula/pic2video/internal/infra/nvenc"
)

func FormatSummary(profile, res, requested, effective, output string, elapsed float64, processed int, warnings []string, nvencAvailable bool) string {
	report := nvenc.BuildReport(requested, effective, nvencAvailable)
	return fmt.Sprintf("status=success profile=%s resolution=%s %s processed=%d elapsed=%.3fs output=%s warnings=%d", profile, res, report, processed, elapsed, output, len(warnings))
}
