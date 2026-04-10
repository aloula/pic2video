package pipeline

import "fmt"

func BuildFramingFilter(width, height int) string {
	return fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", width, height, width, height)
}

func BuildMotionFilter(effect string, width, height int, imageDur float64) string {
	if effect == "" || effect == "static" {
		return ""
	}

	motionFPS := 30
	frames := int(imageDur * float64(motionFPS))
	if frames < 2 {
		frames = 2
	}

	maxDim := width
	if height > maxDim {
		maxDim = height
	}
	resolutionScale := float64(maxDim) / 1080.0
	if resolutionScale < 1.0 {
		resolutionScale = 1.0
	}

	zoomStep := 0.0009
	zoomMax := 1.08
	panFactor := 0.03

	switch effect {
	case "kenburns-medium":
		zoomStep = 0.0015
		zoomMax = 1.15
		panFactor = 0.06
	case "kenburns-high":
		zoomStep = 0.0022
		zoomMax = 1.22
		panFactor = 0.10
	case "kenburns-low":
		// Keep defaults.
	}

	// Keep movement smooth and high quality even with higher processing time.
	zoomStep = zoomStep / resolutionScale
	panScale := panFactor * resolutionScale
	if panScale > 1.0 {
		panScale = 1.0
	}
	prepW := int(float64(width)*1.25 + 0.5)
	prepH := int(float64(height)*1.25 + 0.5)

	return fmt.Sprintf(
		"scale=w=%d:h=%d:force_original_aspect_ratio=increase:flags=lanczos,crop=%d:%d,zoompan=z='min(zoom+%.6f,%.2f)':x='(iw-iw/zoom)*%.4f*(on/%d)':y='(ih-ih/zoom)*%.4f*(on/%d)':d=%d:fps=%d:s=%dx%d,format=yuv420p,setsar=1",
		prepW, prepH,
		prepW, prepH,
		zoomStep, zoomMax,
		panScale,
		frames-1,
		panScale,
		frames-1,
		frames,
		motionFPS,
		width, height,
	)
}
