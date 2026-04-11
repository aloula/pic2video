package pipeline

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/loula/pic2video/internal/domain/media"
)

func BuildXFadeGraph(numAssets int, imageDur, transitionDur float64) string {
	return BuildXFadeGraphWithEffect(numAssets, imageDur, transitionDur, "static", 0, 0)
}

func BuildXFadeGraphWithEffect(numAssets int, imageDur, transitionDur float64, effect string, width, height int) string {
	return BuildXFadeGraphForAssetsWithEffect(nil, numAssets, imageDur, transitionDur, effect, width, height, 60)
}

func BuildXFadeGraphForAssetsWithEffect(assets []media.Asset, numAssets int, imageDur, transitionDur float64, effect string, width, height int, outputFPS int) string {
	if numAssets < 2 {
		return ""
	}
	if outputFPS <= 0 {
		outputFPS = 60
	}
	graph := ""
	startDir := 0
	dirStep := 1
	if effect != "" && effect != "static" {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		startDir = r.Intn(4)
		if r.Intn(2) == 0 {
			dirStep = 3
		}
	}
	for i := 0; i < numAssets; i++ {
		effectForAsset := effect
		asset := media.Asset{}
		if i < len(assets) {
			asset = assets[i]
		}
		if asset.MediaType == media.MediaTypeVideo {
			effectForAsset = "static"
		}
		dir := (startDir + i*dirStep) % 4
		motion := buildMotionFilterWithDirection(effectForAsset, width, height, imageDur, dir)
		if motion == "" {
			framing := BuildFramingFilter(width, height)
			if asset.MediaType == media.MediaTypeVideo {
				framingChain := framing
				if rotFilter := BuildRotationFilter(asset.Rotation); rotFilter != "" {
					framingChain = rotFilter + "," + framing
				}
				graph += fmt.Sprintf("[%d:v]%s,fps=%d,trim=duration=%.3f,setpts=PTS-STARTPTS,format=yuv420p,setsar=1[v%d];", i, framingChain, outputFPS, imageDur, i)
			} else {
				graph += fmt.Sprintf("[%d:v]%s,fps=%d,format=yuv420p,setsar=1[v%d];", i, framing, outputFPS, i)
			}
		} else {
			graph += fmt.Sprintf("[%d:v]%s,fps=%d[v%d];", i, motion, outputFPS, i)
		}
	}
	offset := imageDur - transitionDur
	graph += fmt.Sprintf("[v0][v1]xfade=transition=fade:duration=%.3f:offset=%.3f[x1];", transitionDur, offset)
	for i := 2; i < numAssets; i++ {
		offset += imageDur - transitionDur
		graph += fmt.Sprintf("[x%d][v%d]xfade=transition=fade:duration=%.3f:offset=%.3f[x%d];", i-1, i, transitionDur, offset, i)
	}
	graph += fmt.Sprintf("[x%d]copy[vlast]", numAssets-1)
	return graph
}
