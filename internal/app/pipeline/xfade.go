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
		slotDur := assetSlotDuration(asset, imageDur)
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
				shortVideoFade := buildShortVideoFadeFilter(asset, slotDur)
				padFilter := buildVideoPadFilter(asset, slotDur)
				// Keep timeline offsets stable for mixed-media sequences by extending
				// short video clips up to the slot duration.
				graph += fmt.Sprintf("[%d:v]%s,fps=%d,%s,trim=duration=%.3f,setpts=PTS-STARTPTS%s,format=yuv420p,setsar=1[v%d];", i, framingChain, outputFPS, padFilter, slotDur, shortVideoFade, i)
			} else {
				graph += fmt.Sprintf("[%d:v]%s,fps=%d,format=yuv420p,setsar=1[v%d];", i, framing, outputFPS, i)
			}
		} else {
			graph += fmt.Sprintf("[%d:v]%s,fps=%d[v%d];", i, motion, outputFPS, i)
		}
	}
	firstDur := imageDur
	if len(assets) > 0 {
		firstDur = assetSlotDuration(assets[0], imageDur)
	}
	firstTrans := transitionDurationForEdge(assetAt(assets, 0), transitionDur)
	offset := firstDur - firstTrans
	if offset < 0 {
		offset = 0
	}
	graph += fmt.Sprintf("[v0][v1]xfade=transition=fade:duration=%.3f:offset=%.3f[x1];", firstTrans, offset)
	for i := 2; i < numAssets; i++ {
		prevDur := imageDur
		if i-1 < len(assets) {
			prevDur = assetSlotDuration(assets[i-1], imageDur)
		}
		transDur := transitionDurationForEdge(assetAt(assets, i-1), transitionDur)
		offset += prevDur - transDur
		if offset < 0 {
			offset = 0
		}
		graph += fmt.Sprintf("[x%d][v%d]xfade=transition=fade:duration=%.3f:offset=%.3f[x%d];", i-1, i, transDur, offset, i)
	}
	graph += fmt.Sprintf("[x%d]copy[vlast]", numAssets-1)
	return graph
}

func assetSlotDuration(asset media.Asset, imageDur float64) float64 {
	if asset.MediaType == media.MediaTypeVideo && asset.DurationSec > 0 {
		return asset.DurationSec
	}
	return imageDur
}

func buildShortVideoFadeFilter(asset media.Asset, imageDur float64) string {
	if asset.MediaType != media.MediaTypeVideo {
		return ""
	}
	if asset.DurationSec <= 0 || asset.DurationSec >= imageDur {
		return ""
	}
	fadeDur := 0.25
	if asset.DurationSec < 4 {
		fadeDur = 1.0 / 3.0
	}
	if imageDur/2 < fadeDur {
		fadeDur = imageDur / 2
	}
	if fadeDur <= 0 {
		return ""
	}
	// Fade near source end for all short clips so users don't see a frozen tail.
	fadeOutStart := asset.DurationSec - fadeDur
	if fadeOutStart < 0 {
		fadeOutStart = 0
	}
	return fmt.Sprintf(",fade=t=in:st=0:d=%.3f,fade=t=out:st=%.3f:d=%.3f", fadeDur, fadeOutStart, fadeDur)
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

func assetAt(assets []media.Asset, idx int) media.Asset {
	if idx >= 0 && idx < len(assets) {
		return assets[idx]
	}
	return media.Asset{}
}

func buildVideoPadFilter(asset media.Asset, slotDur float64) string {
	if asset.MediaType == media.MediaTypeVideo && asset.DurationSec > 0 && asset.DurationSec < slotDur {
		// Use black padding for short videos to avoid a frozen last-frame hold.
		return fmt.Sprintf("tpad=stop_mode=add:stop_duration=%.3f:color=black", slotDur)
	}
	return fmt.Sprintf("tpad=stop_mode=clone:stop_duration=%.3f", slotDur)
}
