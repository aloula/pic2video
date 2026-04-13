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
	if !strings.Contains(joined, "fade=t=in:st=0:d=0.500") {
		t.Fatalf("expected global fade-in filter in static mode, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=out:st=8.500:d=0.500") {
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
	if !strings.Contains(joined, "fade=t=in:st=0:d=0.500") {
		t.Fatalf("expected global fade-in filter in kenburns mode, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=out:st=8.500:d=0.500") {
		t.Fatalf("expected global fade-out filter in kenburns mode, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsWithAudioMapsAout(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}, {Path: "c.jpg"}}
	audio := []string{"ambient_a.mp3", "ambient_b.mp3"}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, audio, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "-stream_loop -1") {
		t.Fatalf("did not expect per-input mp3 loop; tracks should play sequentially, got: %s", joined)
	}
	if !strings.Contains(joined, "[3:a][4:a]concat=n=2:v=0:a=1") {
		t.Fatalf("expected concatenated ordered audio labels, got: %s", joined)
	}
	if !strings.Contains(joined, "apad=whole_dur=13.000,atrim=duration=13.000") {
		t.Fatalf("expected padded+bounded audio duration to match slideshow timeline, got: %s", joined)
	}
	if !strings.Contains(joined, "-map [aout]") {
		t.Fatalf("expected mapped audio output, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsWithVideoAudioSourceMapsVideoTrack(t *testing.T) {
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo, HasAudio: true},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPSAndSource("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60, "video")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "[1:a]atrim=duration=9.000,asetpts=N/SR/TB[vsrc]") {
		t.Fatalf("expected video audio source track from video input, got: %s", joined)
	}
	if !strings.Contains(joined, "-map [aout]") {
		t.Fatalf("expected mapped audio output for video audio source, got: %s", joined)
	}
}

// TestBuildRenderCommandArgsShortVideoAudioSourcePadded verifies that when
// --audio-source video is used alongside images, and the video clips are shorter
// than imageDur, apad is inserted before atrim so that -shortest does not
// truncate the output before the image segments are visible.
// Both videos are clamped to imageDur=5s (2 videos + 3 images = 5 slots,
// 4 transitions → total = 25-4 = 21s). Video audio concat = 5.5s which is much
// shorter; without apad, -shortest would end the output at ~5.5s.
func TestBuildRenderCommandArgsShortVideoAudioSourcePadded(t *testing.T) {
	assets := []media.Asset{
		{Path: "v1.mp4", MediaType: media.MediaTypeVideo, DurationSec: 2.5, HasAudio: true},
		{Path: "v2.mp4", MediaType: media.MediaTypeVideo, DurationSec: 3.0, HasAudio: true},
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "b.jpg", MediaType: media.MediaTypeImage},
		{Path: "c.jpg", MediaType: media.MediaTypeImage},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPSAndSource("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60, "video")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "apad=whole_dur=") {
		t.Fatalf("expected apad filter to extend short video audio to full timeline duration, got: %s", joined)
	}
	if !strings.Contains(joined, "atrim=duration=17.833") {
		t.Fatalf("expected final audio trim to updated slideshow duration (17.833s), got: %s", joined)
	}
	if !strings.Contains(joined, "-map [aout]") {
		t.Fatalf("expected audio output mapped, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsWithMixAudioSourceUsesAmix(t *testing.T) {
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo, HasAudio: true},
	}
	audio := []string{"ambient_a.mp3"}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPSAndSource("out.mp4", assets, audio, "static", 5, 1, 1920, 1080, "cpu", 60, "mix")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "amix=inputs=2") {
		t.Fatalf("expected amix chain for mixed audio source, got: %s", joined)
	}
	if !strings.Contains(joined, "-map [aout]") {
		t.Fatalf("expected mapped audio output for mixed audio source, got: %s", joined)
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
	if !strings.Contains(joined, "-an") {
		t.Fatalf("expected -an to mute output when no audio source is available, got: %s", joined)
	}
}

// TestBuildRenderCommandArgsMp3SourceWithVideoAudioMutes verifies that when
// --audio-source mp3 (the default) is selected but there are no MP3 files,
// video audio from input clips is NOT included. The -an flag must be present.
func TestBuildRenderCommandArgsMp3SourceWithVideoAudioMutes(t *testing.T) {
	assets := []media.Asset{
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo, HasAudio: true},
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
	}
	// no audioAssets, source = "mp3" (default)
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPSAndSource("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60, "mp3")
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "[aout]") {
		t.Fatalf("video audio must not be mapped when --audio-source mp3 has no MP3 files, got: %s", joined)
	}
	if !strings.Contains(joined, "-an") {
		t.Fatalf("expected -an to silence video audio when --audio-source mp3 has no MP3 files, got: %s", joined)
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
	if !strings.Contains(joined, "tpad=stop_mode=clone") {
		t.Fatalf("expected video slots to be padded for mixed-media timeline stability, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsLongVideoUsesFullDuration(t *testing.T) {
	assets := []media.Asset{
		{Path: "a.jpg", MediaType: media.MediaTypeImage},
		{Path: "clip.mp4", MediaType: media.MediaTypeVideo, DurationSec: 12},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu")
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "-t 5.000 -i clip.mp4") {
		t.Fatalf("did not expect long video input to be truncated to image duration, got: %s", joined)
	}
	if !strings.Contains(joined, "trim=duration=12.000") {
		t.Fatalf("expected long video slot to preserve full clip duration, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=out:st=15.500:d=0.500") {
		t.Fatalf("expected global fade-out timing based on full mixed timeline duration, got: %s", joined)
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

func TestBuildRenderCommandArgsCPUQualityTuning(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg", MediaType: media.MediaTypeImage}, {Path: "b.jpg", MediaType: media.MediaTypeImage}}
	args := ffmpeg.BuildRenderCommandArgsWithEffect("out.mp4", assets, "static", 5, 1, 1920, 1080, "libx264")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-b:v 6M") || !strings.Contains(joined, "-maxrate 8M") || !strings.Contains(joined, "-bufsize 4M") {
		t.Fatalf("expected improved FHD bitrate tier, got: %s", joined)
	}
	if !strings.Contains(joined, "-preset slower") || !strings.Contains(joined, "-crf 16") {
		t.Fatalf("expected libx264 quality controls (preset+crf), got: %s", joined)
	}
}

func TestBuildRenderCommandArgsNVENCQualityTuning(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg", MediaType: media.MediaTypeImage}, {Path: "b.jpg", MediaType: media.MediaTypeImage}}
	args := ffmpeg.BuildRenderCommandArgsWithEffect("out.mp4", assets, "static", 5, 1, 3840, 2160, "h264_nvenc")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-b:v 20M") || !strings.Contains(joined, "-maxrate 24M") || !strings.Contains(joined, "-bufsize 12M") {
		t.Fatalf("expected improved UHD bitrate tier, got: %s", joined)
	}
	if !strings.Contains(joined, "-preset p6") || !strings.Contains(joined, "-rc vbr_hq") || !strings.Contains(joined, "-cq 17") {
		t.Fatalf("expected nvenc quality controls (preset+rc+cq), got: %s", joined)
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

func TestBuildRenderCommandArgsShortVideosUnder4SecondsUseOneThirdXfade(t *testing.T) {
	assets := []media.Asset{
		{Path: "first_short.mov", MediaType: media.MediaTypeVideo, DurationSec: 2.0},
		{Path: "middle.jpg", MediaType: media.MediaTypeImage},
		{Path: "last_short.mov", MediaType: media.MediaTypeVideo, DurationSec: 1.5},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "xfade=transition=fade:duration=0.333:offset=1.667") {
		t.Fatalf("expected first edge to use 1/3s xfade for <4s video with offset 1.667, got: %s", joined)
	}
	if !strings.Contains(joined, "xfade=transition=fade:duration=1.000:offset=5.667") {
		t.Fatalf("expected second edge to keep configured transition duration (1.000s), got: %s", joined)
	}
	if strings.Contains(joined, "stop_mode=add") {
		t.Fatalf("did not expect black padding for short videos after xfade timing fix, got: %s", joined)
	}
	if !strings.Contains(joined, "fade=t=in:st=0:d=0.500") {
		t.Fatalf("expected global output fade-in to remain 0.5s, got: %s", joined)
	}
}

func TestBuildRenderCommandArgsShortVideoAtLeast4SecondsKeepsConfiguredTransition(t *testing.T) {
	assets := []media.Asset{
		{Path: "first_short.mov", MediaType: media.MediaTypeVideo, DurationSec: 4.0},
		{Path: "middle.jpg", MediaType: media.MediaTypeImage},
	}
	args := ffmpeg.BuildRenderCommandArgsWithEffectAndAudioAndFPS("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", 60)
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "xfade=transition=fade:duration=1.000:offset=3.000") {
		t.Fatalf("expected configured 1.000s xfade for 4s video with offset 3.000, got: %s", joined)
	}
	if strings.Contains(joined, "fade=t=in:st=0:d=0.250") {
		t.Fatalf("did not expect per-clip 0.250 short-video fade-in filter once xfade timing handles transitions, got: %s", joined)
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
