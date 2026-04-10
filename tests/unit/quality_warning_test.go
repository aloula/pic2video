package unit

import (
	"testing"

	"github.com/loula/pic2video/internal/app/renderjob"
	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/domain/profile"
)

func TestEvaluateQualityWarnings(t *testing.T) {
	p, _ := profile.FromName("uhd")
	job := renderjob.RenderJob{Profile: p}
	warnings := renderjob.EvaluateQualityWarnings(job, []media.Asset{{Path: "a.jpg", Width: 1000, Height: 800}})
	if len(warnings) == 0 {
		t.Fatal("expected quality warning")
	}
}
