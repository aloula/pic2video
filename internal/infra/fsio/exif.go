package fsio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ExifData holds extracted EXIF metadata
type ExifData struct {
	CreateDate time.Time
}

// ExtractExif extracts EXIF metadata from an image file using ffprobe
func ExtractExif(imagePath, ffprobeBin string) (*ExifData, error) {
	// FFprobe command to extract image metadata including EXIF
	cmd := exec.Command(ffprobeBin,
		"-v", "error",
		"-print_format", "json",
		"-show_entries", "format_tags=creation_time,com.apple.quicktime.creationtime:format_tags=com.appleqtime.creationtime",
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
	}

	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse exif json: %w", err)
	}

	exif := &ExifData{}

	// Try multiple EXIF date fields (in priority order)
	dateFormats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	// Try creation_time first (most common in modern cameras/phones)
	if creationTime, ok := result.Format.Tags["creation_time"]; ok {
		for _, format := range dateFormats {
			if t, err := time.Parse(format, creationTime); err == nil {
				exif.CreateDate = t
				return exif, nil
			}
		}
	}

	// If creation_time not found or couldn't parse, fall back to file modification time
	fileInfo, err := os.Stat(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	exif.CreateDate = fileInfo.ModTime()

	return exif, nil
}
