package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/app/pipeline"
)

func TestBuildFramingFilter(t *testing.T) {
	f := pipeline.BuildFramingFilter(1920, 1080)
	if !strings.Contains(f, "pad=1920:1080") {
		t.Fatalf("unexpected framing filter: %s", f)
	}
}

func TestBuildMotionFilterStatic(t *testing.T) {
	f := pipeline.BuildMotionFilter("static", 1920, 1080, 5)
	if f != "" {
		t.Fatalf("expected empty static motion filter, got: %s", f)
	}
}

func TestBuildMotionFilterKenBurnsQualityDirectivesAllModes(t *testing.T) {
	for _, mode := range []string{"kenburns-low", "kenburns-medium", "kenburns-high"} {
		t.Run(mode, func(t *testing.T) {
			f := pipeline.BuildMotionFilter(mode, 1920, 1080, 5)
			if !strings.Contains(f, "zoompan=") {
				t.Fatalf("expected zoompan in motion filter for %s: %s", mode, f)
			}
			if !strings.Contains(f, "flags=lanczos") {
				t.Fatalf("expected lanczos scaling in motion filter for %s: %s", mode, f)
			}
			if !strings.Contains(f, "fps=30") {
				t.Fatalf("expected fps=30 smooth motion in filter for %s: %s", mode, f)
			}
			if !strings.Contains(f, "s=1920x1080") {
				t.Fatalf("expected effect output resolution to match profile for %s: %s", mode, f)
			}
			if strings.Contains(f, "crop=") {
				t.Fatalf("ken burns filter must not crop images; use pad instead for %s: %s", mode, f)
			}
			if !strings.Contains(f, "pad=") {
				t.Fatalf("ken burns filter must pad to prevent cropping for %s: %s", mode, f)
			}
		})
	}
}

func TestBuildMotionFilterVariesByAssetIndex(t *testing.T) {
	f0 := pipeline.BuildMotionFilterForAsset("kenburns-medium", 1920, 1080, 5, 0)
	f1 := pipeline.BuildMotionFilterForAsset("kenburns-medium", 1920, 1080, 5, 1)
	f2 := pipeline.BuildMotionFilterForAsset("kenburns-medium", 1920, 1080, 5, 2)
	if f0 == f1 && f1 == f2 {
		t.Fatalf("expected per-asset motion variation, got identical filters")
	}
}

func TestBuildMotionFilterResolutionAware(t *testing.T) {
	fhd := pipeline.BuildMotionFilter("kenburns-high", 1920, 1080, 5)
	uhd := pipeline.BuildMotionFilter("kenburns-high", 3840, 2160, 5)
	if fhd == uhd {
		t.Fatalf("expected different filters for different resolutions")
	}
}

func TestBuildRotationFilter(t *testing.T) {
	cases := []struct {
		degrees int
		want    string
	}{
		{0, ""},
		{-90, "transpose=clock"},
		{270, "transpose=clock"},
		{90, "transpose=cclock"},
		{-270, "transpose=cclock"},
		{180, "hflip,vflip"},
		{-180, "hflip,vflip"},
		{45, ""},
	}
	for _, c := range cases {
		got := pipeline.BuildRotationFilter(c.degrees)
		if got != c.want {
			t.Fatalf("BuildRotationFilter(%d) = %q, want %q", c.degrees, got, c.want)
		}
	}
}
