package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/pipeline"
)

func TestBuildXFadeGraph(t *testing.T) {
	g := pipeline.BuildXFadeGraph(3, 4, 1)
	if !strings.Contains(g, "offset=3.000") {
		t.Fatalf("expected first offset in graph: %s", g)
	}
	if !strings.Contains(g, "offset=6.000") {
		t.Fatalf("expected second offset in graph: %s", g)
	}
}
