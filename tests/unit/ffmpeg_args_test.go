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

func TestBuildRenderCommandArgsOverlayEnabledAddsDrawtext(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio(
		"out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu",
		ffmpeg.OverlayOptions{Enabled: true, FontSize: 42, FooterOffsetPx: 20, BoxAlpha: 0.4, Lines: []string{"A", "B"}},
	)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "drawtext=") {
		t.Fatalf("expected drawtext filter when overlay enabled, got: %s", joined)
	}
	if !strings.Contains(joined, "fontsize=42") {
		t.Fatalf("expected overlay font size in drawtext args, got: %s", joined)
	}
	if !strings.Contains(joined, "y=h-th-30") {
		t.Fatalf("expected 30px footer offset in drawtext args for FHD, got: %s", joined)
	}
	if !strings.Contains(joined, "fontcolor=white") {
		t.Fatalf("expected white overlay text color, got: %s", joined)
	}
	if !strings.Contains(joined, "boxcolor=black@0.40") {
		t.Fatalf("expected semi-transparent box color in drawtext args, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsOverlayDisabledOmitsDrawtext(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio(
		"out.mp4", assets, nil, "static", 5, 1, 3840, 2160, "cpu",
		ffmpeg.OverlayOptions{Enabled: false, FontSize: 60, FooterOffsetPx: 20, BoxAlpha: 0.4, Lines: []string{"A", "B"}},
	)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "drawtext=") {
		t.Fatalf("did not expect drawtext filter when overlay disabled, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsMixedMediaUsesVideoInputWithoutLoop(t *testing.T) {
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if strings.Count(joined, "-loop 1") != 1 {
		t.Fatalf("expected only image inputs to use -loop, got: %s", joined)
	}
	if !strings.Contains(joined, "clip.mp4") {
		t.Fatalf("expected video input to be present, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsFPSPropagation(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg", MediaType: media.MediaTypeImage}, {Path: "b.jpg", MediaType: media.MediaTypeImage}}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 30)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-r 30") {
		t.Fatalf("expected propagated output fps in ffmpeg args, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsMixedMediaKenBurnsSkipsVideoMotion(t *testing.T) {
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "kenburns-medium", 5, 1, 1920, 1080, "cpu", 60)
	joined := strings.Join(args, " ")
	if c := strings.Count(joined, "zoompan="); c != 1 {
		t.Fatalf("expected only image stream to use zoompan, count=%d args=%s", c, joined)
	}
	if !strings.Contains(joined, "fps=60") || !strings.Contains(joined, "trim=duration=5.000") || !strings.Contains(joined, "setpts=PTS-STARTPTS") {
		t.Fatalf("expected video stream normalization before xfade, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsMixedMediaKenBurnsTimbaseNormalized(t *testing.T) {
	// xfade requires identical timebases on both inputs. Ken Burns images produce
	// streams at the zoompan fps (30), while videos are normalized to outputFPS
	// (e.g. 60). Without explicit fps normalization after each asset filter chain
	// the render fails with "First input link main timebase (1/30) do not match
	// the corresponding second input link xfade timebase (1/60)".
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo},
		{Path: "b.jpg", MediaType: media.MediaTypeImage},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "kenburns-low", 5, 1, 1920, 1080, "cpu", 60)
	joined := strings.Join(args, " ")
	// Both image streams (zoompan) and the video stream must produce fps=60 before xfade.
	if c := strings.Count(joined, "fps=60"); c < 3 {
		t.Fatalf("expected at least 3 fps=60 normalization points (2 images + 1 video) before xfade, got %d in: %s", c, joined)
	}
}

func TestBuildRenderCommandArgsStaticPerAssetFramingUsesPad(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg", MediaType: media.MediaTypeImage}, {Path: "b.jpg", MediaType: media.MediaTypeImage}}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if c := strings.Count(joined, "pad=1920:1080"); c < 2 {
		t.Fatalf("expected per-asset fit+pad framing for static mode, got args=%s", joined)
	}
}

func TestBuildRenderCommandArgsKenBurnsNoCropUsesPad(t *testing.T) {
	// Portrait images (and any non-16:9 image) must not be cropped by the Ken Burns
	// filter. Instead the filter must scale-to-fit and pad with black bars.
	assets := []media.Asset{{Path: "portrait.jpg", MediaType: media.MediaTypeImage}, {Path: "landscape.jpg", MediaType: media.MediaTypeImage}}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "kenburns-low", 5, 1, 3840, 2160, "cpu", 60)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, ",crop=") {
		t.Fatalf("ken burns filter must not crop images; use pad instead, got: %s", joined)
	}
	if !strings.Contains(joined, "pad=") {
		t.Fatalf("ken burns filter must include pad= to add black bars, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsVideoPortraitRotationApplied(t *testing.T) {
	// iPhone portrait videos carry rotation=-90 in the Display Matrix side data.
	// The filter graph must prepend transpose=clock to correct orientation before
	// the framing/scale step, otherwise the video plays sideways.
	assets := []media.Asset{
		{Path: "photo.jpg", MediaType: media.MediaTypeImage},
		{Path: "portrait.mov", MediaType: media.MediaTypeVideo, Rotation: -90},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "kenburns-low", 5, 1, 3840, 2160, "cpu", 60)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "transpose=clock") {
		t.Fatalf("expected transpose=clock for Rotation=-90 portrait video, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsVideoNoRotationSkipsTranspose(t *testing.T) {
	assets := []media.Asset{
		{Path: "photo.jpg", MediaType: media.MediaTypeImage},
		{Path: "landscape.mp4", MediaType: media.MediaTypeVideo, Rotation: 0},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60)
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "transpose=") {
		t.Fatalf("did not expect transpose filter for zero-rotation video, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsOverlaySkipsVideoSlots(t *testing.T) {
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo},
		{Path: "b.jpg", MediaType: media.MediaTypeImage},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS(
		"out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60,
		ffmpeg.OverlayOptions{
			Enabled:        true,
			FontSize:       42,
			FooterOffsetPx: 30,
			BoxAlpha:       0.4,
			Lines:          []string{"IMG-A", "", "IMG-B"},
		},
	)
	joined := strings.Join(args, " ")
	if c := strings.Count(joined, "drawtext="); c != 2 {
		t.Fatalf("expected drawtext only for image slots (2), got %d in: %s", c, joined)
	}
	if strings.Contains(joined, "Unknown - Unknown - Unknown") {
		t.Fatalf("did not expect unknown fallback line for explicitly skipped video slot, got: %s", joined)
	}
}
