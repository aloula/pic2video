package pipeline

import (
	"fmt"
	"math/rand"
	"time"
)

func BuildXFadeGraph(numAssets int, imageDur, transitionDur float64) string {
	return BuildXFadeGraphWithEffect(numAssets, imageDur, transitionDur, "static", 0, 0)
}

func BuildXFadeGraphWithEffect(numAssets int, imageDur, transitionDur float64, effect string, width, height int) string {
	if numAssets < 2 {
		return ""
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
		dir := (startDir + i*dirStep) % 4
		motion := buildMotionFilterWithDirection(effect, width, height, imageDur, dir)
		if motion == "" {
			graph += fmt.Sprintf("[%d:v]format=yuv420p,setsar=1[v%d];", i, i)
		} else {
			graph += fmt.Sprintf("[%d:v]%s[v%d];", i, motion, i)
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
