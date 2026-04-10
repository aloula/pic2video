package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/infra/ffmpeg"
)

func TestBuildRenderCommandArgsStaticUsesLoopInputs(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	args := ffmpeg.BuildRenderCommandArgsWithEffect("out.mp4", assets, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-loop 1") {
		t.Fatalf("expected static mode to use looped still inputs, got: %s", joined)
	}
	if !strings.Contains(joined, "-t 5.000") {
		t.Fatalf("expected static mode to include per-image duration input arg, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsKenBurnsUsesSingleFrameInputs(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	args := ffmpeg.BuildRenderCommandArgsWithEffect("out.mp4", assets, "kenburns-high", 5, 1, 3840, 2160, "cpu")
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "-loop 1") {
		t.Fatalf("did not expect looped inputs in kenburns mode, got: %s", joined)
	}
	if strings.Contains(joined, "-t 5.000") {
		t.Fatalf("did not expect per-image -t input arg in kenburns mode, got: %s", joined)
	}
	if !strings.Contains(joined, "zoompan=") {
		t.Fatalf("expected kenburns filter graph to include zoompan, got: %s", joined)
	}
}
