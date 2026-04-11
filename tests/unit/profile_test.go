package unit

import (
	"strings"
	"testing"

	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/domain/profile"
	"github.com/loula/pic2video/internal/infra/ffmpeg"
)

func TestProfileFromName(t *testing.T) {
	p, err := profile.FromName("fhd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Width != 1920 || p.Height != 1080 {
		t.Fatalf("unexpected dimensions: %+v", p)
	}
}

func TestProfileFromNameInvalid(t *testing.T) {
	if _, err := profile.FromName("bad"); err == nil {
		t.Fatal("expected error for invalid profile")
	}
}

func TestOverlayFooterOffsetByProfile(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}}
	overlay := ffmpeg.OverlayOptions{Enabled: true, FontSize: 42, FooterOffsetPx: 20, BoxAlpha: 0.4, Lines: []string{"A", "B"}}
	argsFHD := strings.Join(ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, nil, "static", 5, 1, 1920, 1080, "cpu", overlay), " ")
	argsUHD := strings.Join(ffmpeg.BuildRenderCommandArgsWithEffectAndAudio("out.mp4", assets, nil, "static", 5, 1, 3840, 2160, "cpu", overlay), " ")
	if !strings.Contains(argsFHD, "y=h-th-30") {
		t.Fatalf("expected fhd overlay footer offset y=h-th-30, got: %s", argsFHD)
	}
	if !strings.Contains(argsUHD, "y=h-th-60") {
		t.Fatalf("expected uhd overlay footer offset y=h-th-60, got: %s", argsUHD)
	}
}
