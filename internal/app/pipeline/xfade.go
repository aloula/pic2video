package pipeline

import "fmt"

func BuildXFadeGraph(numAssets int, imageDur, transitionDur float64) string {
	if numAssets < 2 {
		return ""
	}
	graph := ""
	for i := 0; i < numAssets; i++ {
		graph += fmt.Sprintf("[%d:v]format=yuv420p,setsar=1[v%d];", i, i)
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
