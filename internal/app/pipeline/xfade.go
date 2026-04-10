package pipeline

import "fmt"

func BuildXFadeGraph(numAssets int, imageDur, transitionDur float64) string {
	return BuildXFadeGraphWithEffect(numAssets, imageDur, transitionDur, "static", 0, 0)
}

func BuildXFadeGraphWithEffect(numAssets int, imageDur, transitionDur float64, effect string, width, height int) string {
	if numAssets < 2 {
		return ""
	}
	graph := ""
	motion := BuildMotionFilter(effect, width, height, imageDur)
	for i := 0; i < numAssets; i++ {
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
