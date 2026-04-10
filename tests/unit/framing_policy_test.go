package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/pipeline"
)

func TestBuildFramingFilter(t *testing.T) {
	f := pipeline.BuildFramingFilter(1920, 1080)
	if !strings.Contains(f, "pad=1920:1080") {
		t.Fatalf("unexpected framing filter: %s", f)
	}
}
