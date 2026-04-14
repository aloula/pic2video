package ffmpeg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/domain/media"
)

type OverlayOptions struct {
	Enabled        bool
	FontSize       int
	FooterOffsetPx int
	BoxAlpha       float64
	Lines          []string
}

func BuildRenderCommandArgs(outputPath string, assets []media.Asset, imageDur, transitionDur float64, width, height int, encoder string) []string {
	return BuildRenderCommandArgsWithEffect(outputPath, assets, "static", imageDur, transitionDur, width, height, encoder)
}

func BuildRenderCommandArgsWithEffect(outputPath string, assets []media.Asset, imageEffect string, imageDur, transitionDur float64, width, height int, encoder string) []string {
	return BuildRenderCommandArgsWithEffectAndAudio(outputPath, assets, nil, imageEffect, imageDur, transitionDur, width, height, encoder, OverlayOptions{})
}

func BuildRenderCommandArgsWithEffectAndAudio(outputPath string, assets []media.Asset, audioAssets []string, imageEffect string, imageDur, transitionDur float64, width, height int, encoder string, overlays ...OverlayOptions) []string {
	return BuildRenderCommandArgsWithEffectAndAudioAndFPS(outputPath, assets, audioAssets, imageEffect, imageDur, transitionDur, width, height, encoder, 60, overlays...)
}

func BuildRenderCommandArgsWithEffectAndAudioAndFPS(outputPath string, assets []media.Asset, audioAssets []string, imageEffect string, imageDur, transitionDur float64, width, height int, encoder string, outputFPS int, overlays ...OverlayOptions) []string {
	return BuildRenderCommandArgsWithEffectAndAudioAndFPSAndSource(outputPath, assets, audioAssets, imageEffect, imageDur, transitionDur, width, height, encoder, outputFPS, "mp3", overlays...)
}

func BuildRenderCommandArgsWithEffectAndAudioAndFPSAndSource(outputPath string, assets []media.Asset, audioAssets []string, imageEffect string, imageDur, transitionDur float64, width, height int, encoder string, outputFPS int, audioSource string, overlays ...OverlayOptions) []string {
	hasVideo := false
	for _, a := range assets {
		if a.MediaType == media.MediaTypeVideo {
			hasVideo = true
			break
		}
	}
	if hasVideo {
		if transitionDur >= imageDur {
			transitionDur = imageDur * 0.8
		}
		if transitionDur <= 0 {
			transitionDur = imageDur * 0.2
		}
	}
	overlay := OverlayOptions{}
	if len(overlays) > 0 {
		overlay = overlays[0]
	}
	if overlay.Enabled {
		overlay.FooterOffsetPx = overlayFooterOffsetForResolution(width, height)
	}
	inputs := []string{}
	useStaticInputs := imageEffect == "" || imageEffect == "static"
	for _, a := range assets {
		if a.MediaType == media.MediaTypeVideo {
			inputs = append(inputs, "-i", a.Path)
			continue
		}
		if useStaticInputs {
			inputs = append(inputs, "-loop", "1", "-t", fmt.Sprintf("%.3f", imageDur), "-i", a.Path)
		} else {
			// For Ken Burns modes, zoompan controls output duration (`d`) and fps.
			// Avoid loop+t inputs here, which would multiply generated frames.
			inputs = append(inputs, "-i", a.Path)
		}
	}
	for _, audioPath := range audioAssets {
		// Keep each MP3 input finite so multiple files play in order.
		// Duration normalization is handled in the audio filter chain.
		inputs = append(inputs, "-i", audioPath)
	}
	graph := pipeline.BuildXFadeGraphForAssetsWithEffect(assets, len(assets), imageDur, transitionDur, imageEffect, width, height, outputFPS)
	if fadeFilter := buildGlobalFadeFilter(assets, imageDur, transitionDur); fadeFilter != "" {
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + fadeFilter + "[vlast]"
	}
	audioFilter := buildAudioFilter(assets, len(audioAssets), imageDur, transitionDur, audioSource)
	if audioFilter != "" {
		graph += ";" + audioFilter
	}
	if overlayFilter := buildOverlayFilter(overlay, assets, imageDur, transitionDur); overlayFilter != "" {
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + overlayFilter + "[vlast]"
	}
	if outputFPS <= 0 {
		outputFPS = 60
	}
	args := []string{"-y"}
	args = append(args, inputs...)
	args = append(args,
		"-filter_complex", graph,
		"-map", "[vlast]",
		"-c:v", encoder,
	)
	args = append(args, bitrateArgs(width, height)...)
	args = append(args, encoderQualityArgs(encoder)...)
	args = append(args,
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		"-r", fmt.Sprintf("%d", outputFPS),
	)
	if audioFilter != "" {
		args = append(args,
			"-map", "[aout]",
			"-c:a", "aac",
		)
	} else {
		// Explicitly suppress all audio so that input video streams with embedded
		// audio are never auto-selected. Audio is only included when the user
		// chooses --audio-source video or --audio-source mix.
		args = append(args, "-an")
	}
	args = append(args,
		"-shortest",
		filepath.Clean(outputPath),
	)
	return args
}

func escapeDrawtextText(v string) string {
	r := strings.NewReplacer(
		`\\`, `\\\\`,
		`:`, `\\:`,
		`'`, `\\'`,
		`%`, `\\%`,
		`,`, `\\,`,
		`[`, `\\[`,
		`]`, `\\]`,
	)
	return r.Replace(v)
}

func buildOverlayFilter(overlay OverlayOptions, assets []media.Asset, imageDur, transitionDur float64) string {
	assetCount := len(assets)
	if !overlay.Enabled || assetCount == 0 {
		return ""
	}
	fontSize := overlay.FontSize
	if fontSize <= 0 {
		fontSize = 42
	}
	offset := overlay.FooterOffsetPx
	if offset <= 0 {
		offset = 30
	}
	alpha := overlay.BoxAlpha
	if alpha <= 0 || alpha >= 1 {
		alpha = 0.4
	}
	total := totalVideoDuration(assets, imageDur, transitionDur)
	if total <= 0 {
		return ""
	}
	filters := make([]string, 0, assetCount)
	cursor := 0.0
	for i := 0; i < assetCount; i++ {
		dur := assetSlotDuration(assets[i], imageDur)
		start := cursor
		end := total
		edgeTrans := transitionDurationForEdge(assets[i], transitionDur)
		if i < assetCount-1 {
			end = start + (dur - edgeTrans)
		}
		if end <= start {
			end = start + dur
		}
		if end > total {
			end = total
		}
		line := "Unknown - Unknown - Unknown - Unknown - Unknown - Unknown"
		if i < len(overlay.Lines) {
			if strings.TrimSpace(overlay.Lines[i]) == "" {
				continue
			}
			line = overlay.Lines[i]
		}
		filters = append(filters,
			fmt.Sprintf("drawtext=text='%s':fontcolor=white:fontsize=%d:box=1:boxcolor=black@%.2f:boxborderw=12:x=(w-text_w)/2:y=h-th-%d:enable='between(t,%.3f,%.3f)'",
				escapeDrawtextText(line), fontSize, alpha, offset, start, end,
			),
		)
		cursor += dur - edgeTrans
	}
	return strings.Join(filters, ",")
}

func overlayFooterOffsetForResolution(width, height int) int {
	if width*height >= 3840*2160 {
		return 60
	}
	return 30
}

func buildGlobalFadeFilter(assets []media.Asset, imageDur, transitionDur float64) string {
	total := totalVideoDuration(assets, imageDur, transitionDur)
	if total <= 0 {
		return ""
	}
	// Keep fade timing consistent regardless of transition settings.
	fadeDur := 0.5
	maxFade := total / 2
	if fadeDur > maxFade {
		fadeDur = maxFade
	}
	if fadeDur <= 0 {
		return ""
	}
	fadeOutStart := total - fadeDur
	return fmt.Sprintf("fade=t=in:st=0:d=%.3f,fade=t=out:st=%.3f:d=%.3f", fadeDur, fadeOutStart, fadeDur)
}

func bitrateArgs(width, height int) []string {
	var target, maxrate, bufsize string
	if width*height >= 3840*2160 {
		target, maxrate, bufsize = "20M", "24M", "12M"
	} else {
		target, maxrate, bufsize = "6M", "8M", "4M"
	}
	return []string{"-b:v", target, "-maxrate", maxrate, "-bufsize", bufsize}
}

func encoderQualityArgs(encoder string) []string {
	e := strings.ToLower(strings.TrimSpace(encoder))
	if strings.Contains(e, "nvenc") {
		return []string{
			// Balanced defaults: materially faster than quality-first p6/vbr_hq
			// while keeping good visual quality for slideshow content.
			"-preset", "p4",
			"-rc", "vbr",
			"-cq", "19",
			"-spatial_aq", "1",
			"-aq-strength", "8",
			"-temporal_aq", "1",
		}
	}
	// libx264 (or cpu aliases in tests): quality-first VBR.
	return []string{"-preset", "slower", "-crf", "16"}
}

func assetSlotDuration(asset media.Asset, imageDur float64) float64 {
	if asset.MediaType == media.MediaTypeVideo && asset.DurationSec > 0 {
		return asset.DurationSec
	}
	return imageDur
}

func transitionDurationForEdge(asset media.Asset, base float64) float64 {
	if base <= 0 {
		base = 1
	}
	if asset.MediaType == media.MediaTypeVideo && asset.DurationSec > 0 && asset.DurationSec < 4 {
		edge := 1.0 / 3.0
		maxEdge := asset.DurationSec * 0.8
		if maxEdge > 0 && edge > maxEdge {
			edge = maxEdge
		}
		if edge > 0 {
			return edge
		}
	}
	return base
}

func totalVideoDuration(assets []media.Asset, imageDur, transitionDur float64) float64 {
	assetCount := len(assets)
	if assetCount <= 0 {
		return 0
	}
	if assetCount == 1 {
		return assetSlotDuration(assets[0], imageDur)
	}
	total := 0.0
	for i := range assets {
		total += assetSlotDuration(assets[i], imageDur)
	}
	for i := 0; i < assetCount-1; i++ {
		total -= transitionDurationForEdge(assets[i], transitionDur)
	}
	return total
}

func buildAudioFilter(assets []media.Asset, audioInputCount int, imageDur, transitionDur float64, source string) string {
	total := totalVideoDuration(assets, imageDur, transitionDur)
	if total <= 0 {
		return ""
	}

	mode := strings.ToLower(strings.TrimSpace(source))
	if mode == "" {
		mode = "mp3"
	}

	videoLabels := make([]string, 0, len(assets))
	for i, a := range assets {
		if a.MediaType == media.MediaTypeVideo && a.HasAudio {
			videoLabels = append(videoLabels, fmt.Sprintf("[%d:a]", i))
		}
	}

	mp3Labels := make([]string, 0, audioInputCount)
	for i := 0; i < audioInputCount; i++ {
		mp3Labels = append(mp3Labels, fmt.Sprintf("[%d:a]", len(assets)+i))
	}

	videoTrack := buildConcatenatedTrackFilter(videoLabels, total, "vsrc")
	mp3Track := buildConcatenatedTrackFilter(mp3Labels, total, "msrc")

	switch mode {
	case "video":
		if videoTrack == "" {
			return ""
		}
		// apad extends short video audio with silence so the total audio duration
		// matches the full slideshow length. Without it, -shortest would truncate
		// the output at the end of the last video clip, dropping image segments.
		return videoTrack + ";[vsrc]apad=whole_dur=" + fmt.Sprintf("%.3f", total) + ",atrim=duration=" + fmt.Sprintf("%.3f", total) + ",asetpts=N/SR/TB" + buildVideoAudioFadeFilter(total) + "[aout]"
	case "mix":
		switch {
		case videoTrack != "" && mp3Track != "":
			return videoTrack + ";" + mp3Track + ";[vsrc][msrc]amix=inputs=2:duration=longest:dropout_transition=2,apad=whole_dur=" + fmt.Sprintf("%.3f", total) + ",atrim=duration=" + fmt.Sprintf("%.3f", total) + ",asetpts=N/SR/TB" + buildVideoAudioFadeFilter(total) + "[aout]"
		case videoTrack != "":
			// Same apad fix as video-only mode.
			return videoTrack + ";[vsrc]apad=whole_dur=" + fmt.Sprintf("%.3f", total) + ",atrim=duration=" + fmt.Sprintf("%.3f", total) + ",asetpts=N/SR/TB" + buildVideoAudioFadeFilter(total) + "[aout]"
		case mp3Track != "":
			return mp3Track + ";[msrc]apad=whole_dur=" + fmt.Sprintf("%.3f", total) + ",atrim=duration=" + fmt.Sprintf("%.3f", total) + ",asetpts=N/SR/TB" + buildVideoAudioFadeFilter(total) + "[aout]"
		default:
			return ""
		}
	default:
		if mp3Track == "" {
			return ""
		}
		return mp3Track + ";[msrc]apad=whole_dur=" + fmt.Sprintf("%.3f", total) + ",atrim=duration=" + fmt.Sprintf("%.3f", total) + ",asetpts=N/SR/TB" + buildVideoAudioFadeFilter(total) + "[aout]"
	}
}

func buildVideoAudioFadeFilter(total float64) string {
	if total <= 0 {
		return ""
	}
	fadeInDur := 1.0
	fadeOutDur := 3.0
	maxFade := total / 2
	if fadeInDur > maxFade {
		fadeInDur = maxFade
	}
	if fadeOutDur > maxFade {
		fadeOutDur = maxFade
	}
	if fadeInDur <= 0 || fadeOutDur <= 0 {
		return ""
	}
	fadeOutStart := total - fadeOutDur
	return fmt.Sprintf(",afade=t=in:st=0:d=%.3f,afade=t=out:st=%.3f:d=%.3f", fadeInDur, fadeOutStart, fadeOutDur)
}

func buildConcatenatedTrackFilter(labels []string, total float64, out string) string {
	if len(labels) == 0 {
		return ""
	}
	if len(labels) == 1 {
		// atrim caps infinite streams (e.g. -stream_loop mp3); for finite streams
		// (video audio) the caller adds apad so this trim still has the right duration.
		return fmt.Sprintf("%satrim=duration=%.3f,asetpts=N/SR/TB[%s]", labels[0], total, out)
	}
	return fmt.Sprintf("%sconcat=n=%d:v=0:a=1[%s]", strings.Join(labels, ""), len(labels), out)
}
