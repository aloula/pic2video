package unit

import (
	"path/filepath"
	"testing"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/domain/media"
)

func TestApplyOrderName(t *testing.T) {
	assets := []media.Asset{{Path: "b.jpg"}, {Path: "a.jpg"}}
	ordered := pipeline.ApplyOrder("name", assets, nil)
	if filepath.Base(ordered[0].Path) != "a.jpg" {
		t.Fatalf("expected a.jpg first, got %s", ordered[0].Path)
	}
}

func TestApplyOrderExplicit(t *testing.T) {
	assets := []media.Asset{{Path: "a.jpg"}, {Path: "b.jpg"}, {Path: "c.jpg"}}
	ordered := pipeline.ApplyOrder("explicit", assets, []string{"c.jpg", "a.jpg", "b.jpg"})
	if filepath.Base(ordered[0].Path) != "c.jpg" {
		t.Fatalf("expected c.jpg first, got %s", ordered[0].Path)
	}
}
