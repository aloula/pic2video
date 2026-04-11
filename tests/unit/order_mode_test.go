package unit

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/loula/pic2video/internal/app/pipeline"
	"github.com/loula/pic2video/internal/domain/media"
	"github.com/loula/pic2video/internal/infra/fsio"
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

func TestExifNormalizationFallbackUnknown(t *testing.T) {
	if got := fsio.NormalizeExifValue(" "); got != "Unknown" {
		t.Fatalf("expected Unknown fallback for empty value, got=%s", got)
	}
	if got := fsio.NormalizeExifValue("Canon R5"); got != "Canon R5" {
		t.Fatalf("expected non-empty exif value unchanged, got=%s", got)
	}
}

func TestFormatCapturedDate(t *testing.T) {
	if got := fsio.FormatCapturedDate(time.Time{}); got != "Unknown" {
		t.Fatalf("expected Unknown for zero date, got=%s", got)
	}
	tm := time.Date(2024, 12, 31, 10, 11, 12, 0, time.UTC)
	if got := fsio.FormatCapturedDate(tm); got != "31/12/2024" {
		t.Fatalf("expected DD/MM/YYYY formatted date, got=%s", got)
	}
}
