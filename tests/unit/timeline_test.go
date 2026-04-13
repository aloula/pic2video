package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/domain/media"
)

func TestBuildXFadeGraph(t *testing.T) {
	g := pipeline.BuildXFadeGraphWithEffect(3, 4, 1, "static", 1920, 1080)
	if !strings.Contains(g, "offset=3.000") {
		t.Fatalf("expected first offset in graph: %s", g)
	}
	if !strings.Contains(g, "offset=6.000") {
		t.Fatalf("expected second offset in graph: %s", g)
	}
}

func TestBuildXFadeGraphMixedMediaTimingStable(t *testing.T) {
	g := pipeline.BuildXFadeGraphWithEffect(4, 5, 1, "static", 1920, 1080)
	for _, want := range []string{"offset=4.000", "offset=8.000", "offset=12.000"} {
		if !strings.Contains(g, want) {
			t.Fatalf("expected %s in mixed-media timeline graph, got: %s", want, g)
		}
	}
}

func TestBuildXFadeGraphLongVideoKeepsFullClipTiming(t *testing.T) {
	assets := []media.Asset{
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo, DurationSec: 12},
		{Path: "photo.jpg", MediaType: media.MediaTypeImage},
	}
	g := pipeline.BuildXFadeGraphForAssetsWithEffect(assets, len(assets), 5, 1, "static", 1920, 1080, 60)
	if !strings.Contains(g, "trim=duration=12.000") {
		t.Fatalf("expected video slot to keep full duration in filter graph, got: %s", g)
	}
	if !strings.Contains(g, "offset=11.000") {
		t.Fatalf("expected xfade offset to use full first video duration, got: %s", g)
	}
}
