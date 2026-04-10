package ffmpeg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/domain/media"
)

func BuildRenderCommandArgs(outputPath string, assets []media.Asset, imageDur, transitionDur float64, width, height int, encoder string) []string {
	inputs := []string{}
	for _, a := range assets {
		inputs = append(inputs, "-loop", "1", "-t", fmt.Sprintf("%.3f", imageDur), "-i", a.Path)
	}
	graph := pipeline.BuildXFadeGraph(len(assets), imageDur, transitionDur)
	framing := pipeline.BuildFramingFilter(width, height)
	if framing != "" {
		graph = strings.Replace(graph, "[vlast]", "[vtmp]", 1) + ";[vtmp]" + framing + "[vlast]"
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
