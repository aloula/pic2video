package media

// Asset represents one media item in the render timeline.
type Asset struct {
	Path               string
	OrderIndex         int
	Width              int
	Height             int
	Format             string
	IsValid            bool
	ValidationWarnings []string
}
