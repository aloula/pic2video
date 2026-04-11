package pipeline

import "fmt"

func BuildFramingFilter(width, height int) string {
	return fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease:flags=lanczos,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", width, height, width, height)
}

func BuildMotionFilter(effect string, width, height int, imageDur float64) string {
	return BuildMotionFilterForAsset(effect, width, height, imageDur, 0)
}

func BuildMotionFilterForAsset(effect string, width, height int, imageDur float64, assetIndex int) string {
	return buildMotionFilterWithDirection(effect, width, height, imageDur, pseudoRandomDirection(assetIndex))
}

func buildMotionFilterWithDirection(effect string, width, height int, imageDur float64, dir int) string {
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

	xExpr := fmt.Sprintf("(iw-iw/zoom)*%.4f*(on/%d)", panScale, frames-1)
	yExpr := fmt.Sprintf("(ih-ih/zoom)*%.4f*(on/%d)", panScale, frames-1)
	if dir == 1 || dir == 3 {
		xExpr = fmt.Sprintf("(iw-iw/zoom)*%.4f*(1-(on/%d))", panScale, frames-1)
	}
	if dir == 2 || dir == 3 {
		yExpr = fmt.Sprintf("(ih-ih/zoom)*%.4f*(1-(on/%d))", panScale, frames-1)
	}

	return fmt.Sprintf(
		"scale=w=%d:h=%d:force_original_aspect_ratio=decrease:flags=lanczos,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,zoompan=z='min(zoom+%.6f,%.2f)':x='%s':y='%s':d=%d:fps=%d:s=%dx%d,format=yuv420p,setsar=1",
		prepW, prepH,
		prepW, prepH,
		zoomStep, zoomMax,
		xExpr,
		yExpr,
		frames,
		motionFPS,
		width, height,
	)
}

func pseudoRandomDirection(assetIndex int) int {
	if assetIndex < 0 {
		assetIndex = -assetIndex
	}
	v := uint32(assetIndex+1)*2654435761 + 1013904223 //nolint:gomnd
	return int(v % 4)
}

// BuildRotationFilter returns an FFmpeg filter chain fragment that corrects
// the display orientation encoded in the video's Display Matrix side data.
// Returns an empty string when no rotation is needed.
func BuildRotationFilter(rotateDegrees int) string {
	switch rotateDegrees {
	case -90, 270:
		return "transpose=clock"
	case 90, -270:
		return "transpose=cclock"
	case 180, -180:
		return "hflip,vflip"
	default:
		return ""
	}
}
