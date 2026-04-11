package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/loula/pic2video/internal/domain/media"
)

type probeResult struct {
	Streams []struct {
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		Codec      string `json:"codec_name"`
		AvgRate    string `json:"avg_frame_rate"`
		RFrameRate string `json:"r_frame_rate"`
		Tags       struct {
			Rotate string `json:"rotate"`
		} `json:"tags"`
		SideDataList []struct {
			SideDataType string  `json:"side_data_type"`
			Rotation     float64 `json:"rotation"`
		} `json:"side_data_list"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

func parseFPS(raw string) float64 {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "0/0" {
		return 0
	}
	if strings.Contains(raw, "/") {
		parts := strings.Split(raw, "/")
		if len(parts) == 2 {
			n, errN := strconv.ParseFloat(parts[0], 64)
			d, errD := strconv.ParseFloat(parts[1], 64)
			if errN == nil && errD == nil && d > 0 {
				return n / d
			}
		}
	}
	f, _ := strconv.ParseFloat(raw, 64)
	return f
}

func ProbeMedia(ffprobeBin, path string) (media.Asset, error) {
	if ffprobeBin == "" {
		ffprobeBin = "ffprobe"
	}
	// -show_streams emits full stream metadata including side_data_list (Display Matrix rotation)
	// and tags (rotate). -show_entries limits the format section to just duration.
	cmd := exec.Command(ffprobeBin, "-v", "error", "-select_streams", "v:0",
		"-show_streams", "-show_entries", "format=duration",
		"-of", "json", path)
	out, err := cmd.Output()
	if err != nil {
		return media.Asset{}, err
	}
	var pr probeResult
	if err := json.Unmarshal(out, &pr); err != nil {
		return media.Asset{}, err
	}
	if len(pr.Streams) == 0 {
		return media.Asset{}, fmt.Errorf("no stream metadata for %s", path)
	}
	s := pr.Streams[0]
	dur, _ := strconv.ParseFloat(strings.TrimSpace(pr.Format.Duration), 64)
	fps := parseFPS(s.AvgRate)
	if fps <= 0 {
		fps = parseFPS(s.RFrameRate)
	}
	// Resolve display rotation from Display Matrix side data (modern Apple/iPhone files)
	// with fallback to the legacy stream "rotate" tag.
	rotation := 0
	for _, sd := range s.SideDataList {
		if sd.SideDataType == "Display Matrix" && sd.Rotation != 0 {
			rotation = int(sd.Rotation)
			break
		}
	}
	if rotation == 0 && strings.TrimSpace(s.Tags.Rotate) != "" {
		if r, err := strconv.Atoi(strings.TrimSpace(s.Tags.Rotate)); err == nil {
			rotation = r
		}
	}
	return media.Asset{
		Path:        path,
		Width:       s.Width,
		Height:      s.Height,
		DurationSec: dur,
		FrameRate:   fps,
		Format:      s.Codec,
		Rotation:    rotation,
		IsValid:     s.Width > 0 && s.Height > 0,
	}, nil
}

func ProbeImage(ffprobeBin, path string) (media.Asset, error) {
	return ProbeMedia(ffprobeBin, path)
}
