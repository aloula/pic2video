package fsio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
		if v := strings.TrimSpace(tags[strings.ToLower(k)]); v != "" {
			return v
		}
	}
	return ""
}

func normalizeTags(tags map[string]string) map[string]string {
	out := make(map[string]string, len(tags))
	for k, v := range tags {
		lk := strings.ToLower(strings.TrimSpace(k))
		if lk == "" {
			continue
		}
		if _, exists := out[lk]; !exists {
			out[lk] = strings.TrimSpace(v)
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
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006:01:02 15:04:05",
		"2006:01:02 15:04:05-07:00",
		"2006:01:02 15:04:05Z07:00",
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
	// FFprobe command to extract image metadata including EXIF
	cmd := exec.Command(ffprobeBin,
		"-v", "error",
		"-print_format", "json",
		"-show_entries", "format_tags:stream_tags",
		imagePath,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// If ffprobe fails, return error but allow fallback to file mod time
		return nil, fmt.Errorf("failed to extract exif from %s: %w", imagePath, err)
	}

	var result struct {
		Format struct {
			Tags map[string]string `json:"tags"`
		} `json:"format"`
		Streams []struct {
			Tags map[string]string `json:"tags"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse exif json: %w", err)
	}

	tags := normalizeTags(result.Format.Tags)
	for _, s := range result.Streams {
		mergeTags(tags, normalizeTags(s.Tags))
	}

	exif := &ExifData{}
	exif.CameraModel = firstTag(tags, "model", "cameramodelname", "make")
	exif.FocalDistance = firstTag(tags, "focallength", "focallenin35mmfilm", "focal_length")
	exif.ShutterSpeed = firstTag(tags, "exposuretime", "shutterspeedvalue", "shutterspeed")
	exif.Aperture = firstTag(tags, "fnumber", "aperturevalue", "aperture")
	exif.ISO = firstTag(tags, "iso", "isospeedratings", "photographicsensitivity")

	dateRaw := firstTag(
		tags,
		"datetimeoriginal",
		"datetime",
		"createdate",
		"creation_time",
		"com.apple.quicktime.creationtime",
		"com.apple.qtime.creationdate",
	)
	if t, ok := parseExifTime(dateRaw); ok {
		exif.CreateDate = t
		return exif, nil
	}
	if t, ok := parseNumericDateUnix(dateRaw); ok {
		exif.CreateDate = t
		return exif, nil
	}

	// If creation_time not found or couldn't parse, fall back to file modification time
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	exif.CreateDate = fileInfo.ModTime()

	return exif, nil
}
