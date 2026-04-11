package unit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/loula/pic2video/internal/infra/fsio"
)

func writeFakeFFprobe(t *testing.T, json string) string {
	t.Helper()
	d := t.TempDir()
	p := filepath.Join(d, "ffprobe")
	script := "#!/bin/sh\ncat <<'EOF'\n" + json + "\nEOF\n"
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestExtractExifReadsStreamTagsAndParsesDateTimeOriginal(t *testing.T) {
	ffprobe := writeFakeFFprobe(t, `{"format":{"tags":{}},"streams":[{"tags":{"Model":"Canon EOS R5","FocalLength":"35mm","ExposureTime":"1/250","FNumber":"2.8","ISO":"400","DateTimeOriginal":"2024:08:15 12:30:45"}}]}`)
	img := filepath.Join(t.TempDir(), "a.jpg")
	if err := os.WriteFile(img, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	exif, err := fsio.ExtractExif(img, ffprobe)
	if err != nil {
		t.Fatal(err)
	}
	if exif.CameraModel != "Canon EOS R5" {
		t.Fatalf("unexpected model: %s", exif.CameraModel)
	}
	if exif.FocalDistance != "35mm" {
		t.Fatalf("unexpected focal distance: %s", exif.FocalDistance)
	}
	if exif.ShutterSpeed != "1/250" {
		t.Fatalf("unexpected shutter speed: %s", exif.ShutterSpeed)
	}
	if exif.Aperture != "2.8" {
		t.Fatalf("unexpected aperture: %s", exif.Aperture)
	}
	if exif.ISO != "400" {
		t.Fatalf("unexpected iso: %s", exif.ISO)
	}
	if got := fsio.FormatCapturedDate(exif.CreateDate); got != "15/08/2024" {
		t.Fatalf("unexpected parsed capture date: %s", got)
	}
}

func TestExtractExifParsesRFC3339Date(t *testing.T) {
	ffprobe := writeFakeFFprobe(t, `{"format":{"tags":{"creation_time":"2024-08-15T12:30:45.123Z"}},"streams":[{"tags":{}}]}`)
	img := filepath.Join(t.TempDir(), "b.jpg")
	if err := os.WriteFile(img, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	exif, err := fsio.ExtractExif(img, ffprobe)
	if err != nil {
		t.Fatal(err)
	}
	if got := fsio.FormatCapturedDate(exif.CreateDate); got != "15/08/2024" {
		t.Fatalf("unexpected parsed RFC3339 capture date: %s", got)
	}
}

func TestExtractExifFallsBackToFileModTimeWhenDateMissing(t *testing.T) {
	ffprobe := writeFakeFFprobe(t, `{"format":{"tags":{}},"streams":[{"tags":{"Model":"Sony A7"}}]}`)
	img := filepath.Join(t.TempDir(), "c.jpg")
	if err := os.WriteFile(img, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	mod := time.Date(2022, 7, 9, 8, 0, 0, 0, time.UTC)
	if err := os.Chtimes(img, mod, mod); err != nil {
		t.Fatal(err)
	}
	exif, err := fsio.ExtractExif(img, ffprobe)
	if err != nil {
		t.Fatal(err)
	}
	if got := fsio.FormatCapturedDate(exif.CreateDate); got != "09/07/2022" {
		t.Fatalf("expected modtime fallback date, got: %s", got)
	}
}
