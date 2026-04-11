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
	overlay := OverlayOptions{}
	if len(overlays) > 0 {
		overlay = overlays[0]
	}
	inputs := []string{}
	useStaticInputs := imageEffect == "" || imageEffect == "static"
	for _, a := range assets {
		if useStaticInputs {
			inputs = append(inputs, "-loop", "1", "-t", fmt.Sprintf("%.3f", imageDur), "-i", a.Path)
		} else {
			// For Ken Burns modes, zoompan controls output duration (`d`) and fps.
			// Avoid loop+t inputs here, which would multiply generated frames.
			inputs = append(inputs, "-i", a.Path)
		}
	}
	for _, audioPath := range audioAssets {
		inputs = append(inputs, "-i", audioPath)
	}
	graph := pipeline.BuildXFadeGraphWithEffect(len(assets), imageDur, transitionDur, imageEffect, width, height)
	if useStaticInputs {
		framing := pipeline.BuildFramingFilter(width, height)
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + framing + "[vlast]"
	}
	if fadeFilter := buildGlobalFadeFilter(len(assets), imageDur, transitionDur); fadeFilter != "" {
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + fadeFilter + "[vlast]"
	}
	if audioFilter := buildAudioFilter(len(assets), len(audioAssets), imageDur, transitionDur); audioFilter != "" {
		graph += ";" + audioFilter
	}
	if overlayFilter := buildOverlayFilter(overlay, len(assets), imageDur, transitionDur); overlayFilter != "" {
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + overlayFilter + "[vlast]"
	}
	args := []string{"-y"}
	args = append(args, inputs...)
	args = append(args,
		"-filter_complex", graph,
		"-map", "[vlast]",
		"-c:v", encoder,
	)
	args = append(args, bitrateArgs(width, height)...)
	args = append(args,
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		"-r", "60",
	)
	if len(audioAssets) > 0 {
		args = append(args,
			"-map", "[aout]",
			"-c:a", "aac",
		)
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

func buildOverlayFilter(overlay OverlayOptions, assetCount int, imageDur, transitionDur float64) string {
	if !overlay.Enabled || assetCount == 0 {
		return ""
	}
	fontSize := overlay.FontSize
	if fontSize <= 0 {
		fontSize = 42
	}
	offset := overlay.FooterOffsetPx
	if offset <= 0 {
		offset = 10
	}
	alpha := overlay.BoxAlpha
	if alpha <= 0 || alpha >= 1 {
		alpha = 0.4
	}
	total := totalVideoDuration(assetCount, imageDur, transitionDur)
	if total <= 0 {
		return ""
	}
	filters := make([]string, 0, assetCount)
	step := imageDur - transitionDur
	if step <= 0 {
		step = imageDur
	}
	for i := 0; i < assetCount; i++ {
		start := float64(i) * step
		end := start + imageDur
		if end > total {
			end = total
		}
		line := "Unknown - Unknown - Unknown - Unknown - Unknown - Unknown"
		if i < len(overlay.Lines) && strings.TrimSpace(overlay.Lines[i]) != "" {
			line = overlay.Lines[i]
		}
		filters = append(filters,
			fmt.Sprintf("drawtext=text='%s':fontcolor=white:fontsize=%d:box=1:boxcolor=black@%.2f:boxborderw=12:x=(w-text_w)/2:y=h-th-%d:enable='between(t,%.3f,%.3f)'",
				escapeDrawtextText(line), fontSize, alpha, offset, start, end,
			),
		)
	}
	return strings.Join(filters, ",")
}

func buildGlobalFadeFilter(assetCount int, imageDur, transitionDur float64) string {
	total := totalVideoDuration(assetCount, imageDur, transitionDur)
	if total <= 0 {
		return ""
	}
	fadeDur := transitionDur
	if fadeDur <= 0 {
		fadeDur = 0.5
	}
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
		target, maxrate, bufsize = "14M", "16M", "8M"
	} else {
		target, maxrate, bufsize = "3500K", "4M", "2M"
	}
	return []string{"-b:v", target, "-maxrate", maxrate, "-bufsize", bufsize}
}

func totalVideoDuration(assetCount int, imageDur, transitionDur float64) float64 {
	if assetCount <= 0 {
		return 0
	}
	if assetCount == 1 {
		return imageDur
	}
	return float64(assetCount)*imageDur - float64(assetCount-1)*transitionDur
}

func buildAudioFilter(videoInputCount, audioInputCount int, imageDur, transitionDur float64) string {
	if audioInputCount <= 0 {
		return ""
	}
	total := totalVideoDuration(videoInputCount, imageDur, transitionDur)
	if total <= 0 {
		return ""
	}
	labels := make([]string, 0, audioInputCount)
	for i := 0; i < audioInputCount; i++ {
		labels = append(labels, fmt.Sprintf("[%d:a]", videoInputCount+i))
	}
	if audioInputCount == 1 {
		return fmt.Sprintf("%satrim=duration=%.3f,asetpts=N/SR/TB[aout]", labels[0], total)
	}
	return fmt.Sprintf("%sconcat=n=%d:v=0:a=1[aud];[aud]atrim=duration=%.3f,asetpts=N/SR/TB[aout]", strings.Join(labels, ""), audioInputCount, total)
}
