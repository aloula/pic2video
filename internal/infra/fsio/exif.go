package fsio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	goexif "github.com/rwcarlsen/goexif/exif"
)

// ExifData holds extracted EXIF metadata
type ExifData struct {
	CameraModel   string
	FocalDistance string
	ShutterSpeed  string
	Aperture      string
	ISO           string
	CreateDate    time.Time
}

func NormalizeExifValue(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "Unknown"
	}
	return v
}

func FormatCapturedDate(t time.Time) string {
	if t.IsZero() {
		return "Unknown"
	}
	return t.Format("02/01/2006")
}

func firstTag(tags map[string]string, keys ...string) string {
	for _, k := range keys {
		if v := sanitizeExifValue(tags[canonicalTagKey(k)]); v != "" {
			return v
		}
	}
	return ""
}

func sanitizeExifValue(v string) string {
	v = strings.Trim(strings.TrimSpace(v), "\x00")
	for len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			v = strings.TrimSpace(v[1 : len(v)-1])
			continue
		}
		break
	}
	return v
}

func canonicalTagKey(k string) string {
	k = strings.ToLower(strings.TrimSpace(k))
	if k == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(k))
	for _, r := range k {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeTags(tags map[string]string) map[string]string {
	out := make(map[string]string, len(tags))
	for k, v := range tags {
		lk := canonicalTagKey(k)
		if lk == "" {
			continue
		}
		if _, exists := out[lk]; !exists {
			out[lk] = sanitizeExifValue(v)
		}
	}
	return out
}

func mergeTags(dst map[string]string, src map[string]string) {
	for k, v := range src {
		if strings.TrimSpace(v) == "" {
			continue
		}
		if _, exists := dst[k]; !exists {
			dst[k] = v
		}
	}
}

func parseExifTime(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000000Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05-0700",
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05-0700",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02",
		"2006:01:02 15:04:05",
		"2006:01:02 15:04:05.000",
		"2006:01:02 15:04:05-07:00",
		"2006:01:02 15:04:05Z07:00",
		"2006:01:02 15:04:05-0700",
		"2006:01:02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, true
		}
	}
	if i := strings.Index(raw, "."); i > 0 {
		trimmed := raw[:i]
		if t, ok := parseExifTime(trimmed); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

func applyNativeExifFallback(imagePath string, exif *ExifData) {
	if exif == nil {
		return
	}
	ext := strings.ToLower(filepath.Ext(imagePath))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".tif" && ext != ".tiff" {
		return
	}

	f, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer f.Close()

	x, err := goexif.Decode(f)
	if err != nil {
		return
	}

	if exif.CameraModel == "" {
		if tag, err := x.Get(goexif.Model); err == nil {
			if model, err := tag.StringVal(); err == nil {
				exif.CameraModel = sanitizeExifValue(model)
			}
		}
	}
	if exif.FocalDistance == "" {
		if tag, err := x.Get(goexif.FocalLength); err == nil {
			exif.FocalDistance = sanitizeExifValue(tag.String())
		}
	}
	if exif.ShutterSpeed == "" {
		if tag, err := x.Get(goexif.ExposureTime); err == nil {
			exif.ShutterSpeed = sanitizeExifValue(tag.String())
		}
	}
	if exif.Aperture == "" {
		if tag, err := x.Get(goexif.FNumber); err == nil {
			exif.Aperture = sanitizeExifValue(tag.String())
		}
	}
	if exif.ISO == "" {
		if tag, err := x.Get(goexif.ISOSpeedRatings); err == nil {
			exif.ISO = sanitizeExifValue(tag.String())
		}
	}
	if exif.CreateDate.IsZero() {
		if t, err := x.DateTime(); err == nil {
			exif.CreateDate = t
		}
	}
}

func parseNumericDateUnix(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	if sec, err := strconv.ParseInt(raw, 10, 64); err == nil && sec > 0 {
		return time.Unix(sec, 0), true
	}
	return time.Time{}, false
}

// ExtractExif extracts EXIF metadata from an image file using ffprobe
func ExtractExif(imagePath, ffprobeBin string) (*ExifData, error) {
	if strings.TrimSpace(ffprobeBin) == "" {
		ffprobeBin = "ffprobe"
	}
	// FFprobe command to extract media metadata including EXIF/video tags.
	// -show_format/-show_streams are required so format_tags/stream_tags entries
	// are actually emitted across different ffprobe builds.
	cmd := exec.Command(ffprobeBin,
		"-v", "error",
		"-show_format",
		"-show_streams",
		"-print_format", "json",
		"-show_entries", "format_tags:stream_tags",
		imagePath,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ffprobeErr := cmd.Run()

	var result struct {
		Format struct {
			Tags map[string]string `json:"tags"`
		} `json:"format"`
		Streams []struct {
			Tags map[string]string `json:"tags"`
		} `json:"streams"`
	}

	if ffprobeErr == nil {
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			return nil, fmt.Errorf("failed to parse exif json: %w", err)
		}
	}

	tags := normalizeTags(result.Format.Tags)
	for _, s := range result.Streams {
		mergeTags(tags, normalizeTags(s.Tags))
	}

	exif := &ExifData{}
	exif.CameraModel = firstTag(tags, "model", "cameramodelname", "com.apple.quicktime.model", "make")
	exif.FocalDistance = firstTag(tags, "focallength", "focallenin35mmfilm", "focallengthin35mmfilm", "focal_length")
	exif.ShutterSpeed = firstTag(tags, "exposuretime", "shutterspeedvalue", "shutterspeed")
	exif.Aperture = firstTag(tags, "fnumber", "aperturevalue", "aperture")
	exif.ISO = firstTag(tags, "iso", "isospeedratings", "isoequivalent", "photographicsensitivity")

	dateRaw := firstTag(
		tags,
		"datetimeoriginal",
		"datetime",
		"date",
		"createdate",
		"creationdate",
		"mediacreatedate",
		"creation_time",
		"com.apple.quicktime.creationtime",
		"com.apple.qtime.creationdate",
		"com.apple.quicktime.creationdate",
	)
	if t, ok := parseExifTime(dateRaw); ok {
		exif.CreateDate = t
	}
	if exif.CreateDate.IsZero() {
		if t, ok := parseNumericDateUnix(dateRaw); ok {
			exif.CreateDate = t
		}
	}

	applyNativeExifFallback(imagePath, exif)

	if exif.CreateDate.IsZero() {
		// If creation date is still unavailable, fall back to file modification time
		fileInfo, err := os.Stat(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file: %w", err)
		}
		exif.CreateDate = fileInfo.ModTime()
	}

	if ffprobeErr != nil && exif.CameraModel == "" && exif.FocalDistance == "" && exif.ShutterSpeed == "" && exif.Aperture == "" && exif.ISO == "" {
		ffprobeStderr := strings.TrimSpace(stderr.String())
		if ffprobeStderr != "" {
			return exif, fmt.Errorf("failed to extract exif from %s: %w (%s)", imagePath, ffprobeErr, ffprobeStderr)
		}
		return exif, fmt.Errorf("failed to extract exif from %s: %w", imagePath, ffprobeErr)
	}

	return exif, nil
}
