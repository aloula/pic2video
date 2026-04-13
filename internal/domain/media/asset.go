package media

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

// Asset represents one media item in the render timeline.
type Asset struct {
	Path               string
	MediaType          MediaType
	OrderIndex         int
	Width              int
	Height             int
	DurationSec        float64
	FrameRate          float64
	HasAudio           bool
	Format             string
	Rotation           int
	IsValid            bool
	ValidationWarnings []string
}
