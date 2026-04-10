package ffmpeg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/domain/media"
)

func BuildRenderCommandArgs(outputPath string, assets []media.Asset, imageDur, transitionDur float64, width, height int, encoder string) []string {
	return BuildRenderCommandArgsWithEffect(outputPath, assets, "static", imageDur, transitionDur, width, height, encoder)
}

func BuildRenderCommandArgsWithEffect(outputPath string, assets []media.Asset, imageEffect string, imageDur, transitionDur float64, width, height int, encoder string) []string {
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
	graph := pipeline.BuildXFadeGraphWithEffect(len(assets), imageDur, transitionDur, imageEffect, width, height)
	if useStaticInputs {
		framing := pipeline.BuildFramingFilter(width, height)
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + framing + "[vlast]"
	}
	if fadeFilter := buildGlobalFadeFilter(len(assets), imageDur, transitionDur); fadeFilter != "" {
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + fadeFilter + "[vlast]"
	}
	args := []string{"-y"}
	args = append(args, inputs...)
	args = append(args,
		"-filter_complex", graph,
		"-map", "[vlast]",
		"-c:v", encoder,
		"-pix_fmt", "yuv420p",
		"-movflags", "+faststart",
		"-r", "30",
		"-shortest",
		filepath.Clean(outputPath),
	)
	return args
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

func totalVideoDuration(assetCount int, imageDur, transitionDur float64) float64 {
	if assetCount <= 0 {
		return 0
	}
	if assetCount == 1 {
		return imageDur
	}
	return float64(assetCount)*imageDur - float64(assetCount-1)*transitionDur
}
