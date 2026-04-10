package pipeline

import "fmt"

func BuildFramingFilter(width, height int) string {
	return fmt.Sprintf("scale=w=%d:h=%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", width, height, width, height)
}
