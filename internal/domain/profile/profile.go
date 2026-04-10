package profile

import "fmt"

// Profile defines a supported output profile.
type Profile struct {
	Name        string
	Width       int
	Height      int
	AspectRatio string
}

const (
FHD = "fhd"
UHD = "uhd"
)

func FromName(name string) (Profile, error) {
	switch name {
	case FHD:
		return Profile{Name: FHD, Width: 1920, Height: 1080, AspectRatio: "16:9"}, nil
	case UHD:
		return Profile{Name: UHD, Width: 3840, Height: 2160, AspectRatio: "16:9"}, nil
	default:
		return Profile{}, fmt.Errorf("invalid profile: %s", name)
	}
}
