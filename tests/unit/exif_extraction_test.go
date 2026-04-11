package unit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestExtractExifReadsQuickTimeVideoTags(t *testing.T) {
	ffprobe := writeFakeFFprobe(t, `{"format":{"tags":{"com.apple.quicktime.model":"iPhone 15 Pro","com.apple.quicktime.creationdate":"2026-01-18T12:40:12+0900"}},"streams":[{"tags":{}}]}`)
	video := filepath.Join(t.TempDir(), "clip.mov")
	if err := os.WriteFile(video, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	exif, err := fsio.ExtractExif(video, ffprobe)
	if err != nil {
		t.Fatal(err)
	}
	if exif.CameraModel != "iPhone 15 Pro" {
		t.Fatalf("unexpected model from quicktime tags: %s", exif.CameraModel)
	}
	if got := fsio.FormatCapturedDate(exif.CreateDate); got != "18/01/2026" {
		t.Fatalf("unexpected parsed quicktime capture date: %s", got)
	}
}

func TestExtractExifIncludesShowFormatAndShowStreamsFlags(t *testing.T) {
	t.Helper()
	d := t.TempDir()
	argsLog := filepath.Join(d, "args.log")
	p := filepath.Join(d, "ffprobe")
	script := fmt.Sprintf(`#!/bin/sh
printf "%%s\n" "$@" > %q
cat <<'EOF'
{"format":{"tags":{"model":"Canon"}},"streams":[{"tags":{}}]}
EOF
`, argsLog)
	if err := os.WriteFile(p, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	img := filepath.Join(d, "a.jpg")
	if err := os.WriteFile(img, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := fsio.ExtractExif(img, p); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(argsLog)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{"-show_format", "-show_streams", "-show_entries", "format_tags:stream_tags"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected ffprobe args to include %s, got:\n%s", want, got)
		}
	}
}
