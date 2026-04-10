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
	if !strings.Contains(joined, "fade=t=in:st=0:d=1.000") {
		t.Fatalf("expected global fade-in filter in static mode, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=out:st=8.000:d=1.000") {
		t.Fatalf("expected global fade-out filter in static mode, got: %s", joined)
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
	if !strings.Contains(joined, "fade=t=in:st=0:d=1.000") {
		t.Fatalf("expected global fade-in filter in kenburns mode, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=out:st=8.000:d=1.000") {
		t.Fatalf("expected global fade-out filter in kenburns mode, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsWithAudioMapsAout(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}, {Path: "c.jpg"}}
	audio := []string{"ambient_a.mp3", "ambient_b.mp3"}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, audio, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "[3:a][4:a]concat=n=2:v=0:a=1") {
		t.Fatalf("expected concatenated ordered audio labels, got: %s", joined)
	}
	if !strings.Contains(joined, "atrim=duration=13.000") {
		t.Fatalf("expected bounded audio duration to match slideshow timeline, got: %s", joined)
	}
	if !strings.Contains(joined, "-map [aout]") {
		t.Fatalf("expected mapped audio output, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsShortDurationFadeClamp(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	args := ffmpeg.BuildRenderCommandArgsWithEffect("out.mp4", assets, "static", 1, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "fade=t=in:st=0:d=0.500") {
		t.Fatalf("expected fade-in clamp to half total duration, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=out:st=0.500:d=0.500") {
		t.Fatalf("expected fade-out clamp timing for short duration, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsNoAudioLeavesVideoOnlyMap(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "[aout]") {
		t.Fatalf("did not expect audio map when no mp3 assets are present, got: %s", joined)
	}
}
