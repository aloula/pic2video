package unit

import (
	"testing"

	"github.com/loula/pic2video/internal/domain/profile"
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
