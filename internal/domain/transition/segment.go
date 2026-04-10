package transition

// Segment represents one transition from asset i to asset i+1.
type Segment struct {
	FromAssetIndex int
	ToAssetIndex   int
	Type           string
	DurationSec    float64
	OffsetSec      float64
}
